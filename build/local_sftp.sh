#!/usr/bin/env bash
#
# Copyright 2025 The Kanister Authors.
# Wrapper for managing a local SFTP server for testing

set -o errexit
set -o nounset
set -o xtrace
set -o pipefail

readonly BASE_DIR=$(dirname ${0})
readonly SFTP_CONTAINER_NAME=${SFTP_CONTAINER_NAME:-"kanister-sftp-test"}
readonly SFTP_IMAGE=${SFTP_IMAGE:-"atmoz/sftp:latest"}
readonly SFTP_PORT=${SFTP_PORT:-"2222"}
readonly SFTP_DATA_DIR=${SFTP_DATA_DIR:-"/tmp/kanister-sftp-test"}
readonly SFTP_SSH_DIR=${SFTP_SSH_DIR:-"/tmp/kanister-sftp-test-ssh"}

# User credentials
readonly SFTP_USER_PASSWORD="sftpuser"
readonly SFTP_USER_PASSWORD_PASS="sftppass123"
readonly SFTP_USER_KEY="sftpkey"

check_docker() {
    if ! command -v docker > /dev/null 2>&1; then
        echo "Error: docker is required but not installed"
        exit 1
    fi
}

generate_ssh_keys() {
    echo "Generating SSH keys for key-based authentication..."
    mkdir -p "${SFTP_SSH_DIR}"
    
    # Generate SSH key for key-based auth user if it doesn't exist
    if [ ! -f "${SFTP_SSH_DIR}/id_rsa" ]; then
        ssh-keygen -t rsa -b 2048 -f "${SFTP_SSH_DIR}/id_rsa" -N "" -C "kanister-sftp-test"
        echo "Generated SSH key pair at ${SFTP_SSH_DIR}/id_rsa"
    else
        echo "SSH key already exists at ${SFTP_SSH_DIR}/id_rsa"
    fi
    
    # Create authorized_keys file for the container
    cat "${SFTP_SSH_DIR}/id_rsa.pub" > "${SFTP_SSH_DIR}/authorized_keys"
    chmod 600 "${SFTP_SSH_DIR}/authorized_keys"
    
    # Create users.conf for atmoz/sftp
    # Format: username:password:uid:gid:upload_dir
    # For key-based auth user, use encrypted password (x means use SSH keys only)
    cat > "${SFTP_SSH_DIR}/users.conf" <<EOF
${SFTP_USER_PASSWORD}:${SFTP_USER_PASSWORD_PASS}:1001:1001
${SFTP_USER_KEY}::1002:1002
EOF
    
    echo "Created users configuration"
}

start_sftp() {
    echo "Starting local SFTP server..."
    check_docker
    
    # Stop any existing container
    if docker ps -a --format '{{.Names}}' | grep -q "^${SFTP_CONTAINER_NAME}$"; then
        echo "Removing existing SFTP container..."
        docker rm -f "${SFTP_CONTAINER_NAME}" || true
    fi
    
    # Create data directory
    mkdir -p "${SFTP_DATA_DIR}"
    chmod 755 "${SFTP_DATA_DIR}"
    
    # Generate SSH keys and user config
    generate_ssh_keys
    
    # Start SFTP container
    echo "Starting SFTP container on port ${SFTP_PORT}..."
    docker run -d \
        --name "${SFTP_CONTAINER_NAME}" \
        -p "${SFTP_PORT}:22" \
        -v "${SFTP_SSH_DIR}/users.conf:/etc/sftp/users.conf:ro" \
        -v "${SFTP_SSH_DIR}/authorized_keys:/home/${SFTP_USER_KEY}/.ssh/keys/id_rsa.pub:ro" \
        "${SFTP_IMAGE}"
    
    # Wait for container to be ready
    echo "Waiting for SFTP server to be ready..."
    local retries=30
    while ! docker exec "${SFTP_CONTAINER_NAME}" pgrep sshd > /dev/null 2>&1; do
        if [[ ${retries} -le 0 ]]; then
            echo "Error: SFTP server failed to start"
            docker logs "${SFTP_CONTAINER_NAME}"
            return 1
        fi
        sleep 1
        retries=$((retries-1))
    done
    
    # Create upload directories with proper permissions
    docker exec "${SFTP_CONTAINER_NAME}" mkdir -p "/home/${SFTP_USER_PASSWORD}/upload"
    docker exec "${SFTP_CONTAINER_NAME}" chown 1001:1001 "/home/${SFTP_USER_PASSWORD}/upload"
    docker exec "${SFTP_CONTAINER_NAME}" mkdir -p "/home/${SFTP_USER_KEY}/upload"
    docker exec "${SFTP_CONTAINER_NAME}" chown 1002:1002 "/home/${SFTP_USER_KEY}/upload"
    
    echo ""
    echo "================================================================"
    echo "SFTP server started successfully!"
    echo "================================================================"
    echo "Container name: ${SFTP_CONTAINER_NAME}"
    echo "Host: localhost"
    echo "Port: ${SFTP_PORT}"
    echo ""
    echo "User 1 (password authentication):"
    echo "  Username: ${SFTP_USER_PASSWORD}"
    echo "  Password: ${SFTP_USER_PASSWORD_PASS}"
    echo "  Home dir: /home/${SFTP_USER_PASSWORD}/upload"
    echo ""
    echo "User 2 (key-based authentication):"
    echo "  Username: ${SFTP_USER_KEY}"
    echo "  Private key: ${SFTP_SSH_DIR}/id_rsa"
    echo "  Home dir: /home/${SFTP_USER_KEY}/upload"
    echo ""
    echo "Test connection with password:"
    echo "  sftp -P ${SFTP_PORT} ${SFTP_USER_PASSWORD}@localhost"
    echo ""
    echo "Test connection with key:"
    echo "  sftp -P ${SFTP_PORT} -i ${SFTP_SSH_DIR}/id_rsa ${SFTP_USER_KEY}@localhost"
    echo "================================================================"
}

stop_sftp() {
    echo "Stopping local SFTP server..."
    check_docker
    
    if docker ps -a --format '{{.Names}}' | grep -q "^${SFTP_CONTAINER_NAME}$"; then
        docker rm -f "${SFTP_CONTAINER_NAME}"
        echo "SFTP container stopped and removed"
    else
        echo "SFTP container not found"
    fi
}

status_sftp() {
    check_docker
    
    if docker ps --format '{{.Names}}' | grep -q "^${SFTP_CONTAINER_NAME}$"; then
        echo "SFTP server is running"
        echo ""
        docker ps --filter "name=${SFTP_CONTAINER_NAME}" --format "table {{.ID}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}"
        echo ""
        echo "Connection details:"
        echo "  Host: localhost"
        echo "  Port: ${SFTP_PORT}"
        echo "  User (password): ${SFTP_USER_PASSWORD} / ${SFTP_USER_PASSWORD_PASS}"
        echo "  User (key): ${SFTP_USER_KEY} (key at ${SFTP_SSH_DIR}/id_rsa)"
        return 0
    else
        echo "SFTP server is not running"
        return 1
    fi
}

logs_sftp() {
    check_docker
    
    if docker ps -a --format '{{.Names}}' | grep -q "^${SFTP_CONTAINER_NAME}$"; then
        docker logs "${SFTP_CONTAINER_NAME}" "$@"
    else
        echo "SFTP container not found"
        return 1
    fi
}

clean_sftp() {
    echo "Cleaning up SFTP data and SSH keys..."
    stop_sftp
    
    if [ -d "${SFTP_DATA_DIR}" ]; then
        rm -rf "${SFTP_DATA_DIR}"
        echo "Removed data directory: ${SFTP_DATA_DIR}"
    fi
    
    if [ -d "${SFTP_SSH_DIR}" ]; then
        rm -rf "${SFTP_SSH_DIR}"
        echo "Removed SSH directory: ${SFTP_SSH_DIR}"
    fi
    
    echo "Cleanup complete"
}

usage() {
    cat <<EOM
Usage: ${0} <operation>
Where operation is one of the following:
  start_sftp  : Start local SFTP server
  stop_sftp   : Stop local SFTP server
  status_sftp : Show SFTP server status
  logs_sftp   : Show SFTP server logs
  clean_sftp  : Stop server and clean up all data
EOM
    exit 1
}

[ ${#@} -gt 0 ] || usage
check_docker
case "${1}" in
        # Alphabetically sorted
        clean_sftp)
            time -p clean_sftp
            ;;
        logs_sftp)
            shift
            logs_sftp "$@"
            ;;
        start_sftp)
            time -p start_sftp
            ;;
        status_sftp)
            status_sftp
            ;;
        stop_sftp)
            time -p stop_sftp
            ;;
        *)
            usage
            exit 1
esac
