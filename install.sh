#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

BINARY_NAME="nullbeacon"
INSTALL_DIR="/usr/local/bin"
DOWNLOAD_URL="https://github.com/RodKast/NullBeacon/releases/latest/download/nullbeacon-linux-amd64"

echo -e "${BOLD}${CYAN}"
echo "  _   _       _ _ ____"
echo " | \ | |_   _| | | __ )  ___  __ _  ___ ___  _ __"
echo " |  \| | | | | | |  _ \ / _ \/ _\` |/ __/ _ \| '_ \\"
echo " | |\  | |_| | | | |_) |  __/ (_| | (_| (_) | | | |"
echo " |_| \_|\__,_|_|_|____/ \___|\__,_|\___\___/|_| |_|"
echo -e "${RESET}"
echo -e "${BOLD}  C2 Framework — github.com/RodKast/NullBeacon${RESET}"
echo ""

if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}[!] This script must be run as root.${RESET}"
    echo -e "    Try: ${BOLD}sudo bash install.sh${RESET}"
    exit 1
fi

echo -e "${CYAN}[*] Downloading NullBeacon...${RESET}"
curl -fsSL -o "/tmp/${BINARY_NAME}" "${DOWNLOAD_URL}"

echo -e "${CYAN}[*] Installing to ${INSTALL_DIR}/${BINARY_NAME}...${RESET}"
mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

echo ""
echo -e "${GREEN}[+] NullBeacon installed successfully!${RESET}"
echo -e "    Run ${BOLD}nullbeacon${RESET} to start the teamserver."
echo ""
