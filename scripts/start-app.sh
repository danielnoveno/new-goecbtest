#!/bin/bash

# ECB Test Start Script
# Digunakan untuk menjalankan aplikasi secara otomatis saat boot.

# Tentukan path aplikasi (sesuaikan jika folder project berbeda)
APP_DIR="/home/pi/V2go-ecbtest"
BINARY="./bin/ecom"

# Masuk ke direktori aplikasi agar file .env dan assets terbaca
cd $APP_DIR || exit

# Pastikan environment display sudah ada (penting untuk GUI via autostart)
export DISPLAY=:0

# Log output ke file untuk debugging
LOG_FILE="$APP_DIR/app.log"
echo "--- App started at $(date) ---" >> $LOG_FILE

# Cek apakah binary ada
if [ -f "$BINARY" ]; then
    echo "Running binary..." >> $LOG_FILE
    $BINARY >> $LOG_FILE 2>&1
else
    echo "Error: Binary $BINARY tidak ditemukan di $APP_DIR" >> $LOG_FILE
    exit 1
fi
