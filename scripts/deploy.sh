#!/bin/bash

# ==========================================
# Script Automasi Deploy Project ke VPS
# ==========================================

# 1. Konfigurasi VPS (UBAH BAGIAN INI SESUAI VPS KAMU)
VPS_USER="root"               # Username SSH (misal: root, ubuntu, debian)
VPS_IP="202.10.34.21"         # IP Address VPS kamu
VPS_PORT="22"                 # Port SSH (default: 22)
REMOTE_DIR="/opt/alibabab"    # Direktori tujuan di VPS

# Warna untuk output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🚀 Memulai proses deploy ke VPS ($VPS_USER@$VPS_IP)...${NC}"

# 2. Membuat direktori di VPS
echo -e "\n${YELLOW}📁 Membuat direktori target di VPS...${NC}"
ssh -p $VPS_PORT $VPS_USER@$VPS_IP "mkdir -p $REMOTE_DIR"
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Gagal terhubung ke VPS. Pastikan IP, Username, dan koneksi internet kamu aman.${NC}"
    exit 1
fi

# 3. Copy file project menggunakan rsync
echo -e "\n${YELLOW}🔼 Mengupload file ke VPS (mengabaikan direktori lokal seperti node_modules dan .git)...${NC}"
rsync -avz --progress -e "ssh -p $VPS_PORT" \
    --exclude='.git' \
    --exclude='node_modules' \
    --exclude='frontend/.next' \
    --exclude='backend/bin' \
    --exclude='backend/tmp' \
    --exclude='postgres_data' \
    --exclude='redis_data' \
    ./ $VPS_USER@$VPS_IP:$REMOTE_DIR/

# 4. Copy custom .env (Mencegah timpa .env jika sudah ada, atau buat dari .env.example)
echo -e "\n${YELLOW}⚙️ Mengecek konfigurasi .env di VPS...${NC}"
ssh -p $VPS_PORT $VPS_USER@$VPS_IP "cd $REMOTE_DIR && if [ ! -f .env ]; then cp .env.example .env; echo '.env baru saja dibuat dari template.'; fi"

# 5. Build & Run Docker Compose di VPS
echo -e "\n${YELLOW}🐳 Build & Start Docker Containers di VPS...${NC}"
ssh -p $VPS_PORT $VPS_USER@$VPS_IP "cd $REMOTE_DIR && docker compose -f docker-compose.yml up -d --build"

echo -e "\n${GREEN}✅ Deploy berhasil diselesaikan!${NC}"
echo -e "Akses websitemu di: ${YELLOW}http://$VPS_IP:3000${NC} (Frontend) dan ${YELLOW}http://$VPS_IP:8080${NC} (Backend)"
echo -e "\nJika butuh edit environment (.env), jalankan:"
echo -e "${YELLOW}ssh -p $VPS_PORT $VPS_USER@$VPS_IP${NC}"
echo -e "${YELLOW}nano $REMOTE_DIR/.env${NC}"
echo -e "${YELLOW}cd $REMOTE_DIR && docker compose restart${NC}"
