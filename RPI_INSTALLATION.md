# Panduan Instalasi dan Menjalankan Aplikasi di Raspberry Pi

Dokumen ini berisi langkah-langkah lengkap untuk menginstal MySQL (MariaDB) dan menjalankan aplikasi **V2go-ecbtest** di Raspberry Pi dari nol.

---

## 1. Persiapan Sistem
Pertama, pastikan sistem Raspberry Pi Anda dalam keadaan terbaru.

```bash
sudo apt update
sudo apt upgrade -y
```

---

## 2. Instalasi dan Konfigurasi Database (MariaDB)
Raspberry Pi umumnya menggunakan MariaDB sebagai alternatif MySQL yang lebih ringan dan kompatibel.

### A. Instal MariaDB Server
```bash
sudo apt install mariadb-server -y
```

### B. Amankan Instalasi
Jalankan perintah ini untuk mengatur password root dan menghapus akses yang tidak perlu:
```bash
sudo mysql_secure_installation
```
*Ikuti petunjuk di layar (pilih 'Y' untuk semua saran, dan atur password root).*

### C. Buat Database dan User Aplikasi
Masuk ke terminal MySQL:
```bash
sudo mysql -u root -p
```

Di dalam shell MySQL/MariaDB, jalankan perintah berikut (sesuaikan dengan `.env` Anda):
```sql
-- Buat database
CREATE DATABASE ecbtest;

-- Buat user dan berikan akses (Gunakan password yang kuat)
-- Catatan: Menggunakan mysql_native_password agar kompatibel dengan library Go
CREATE USER 'ecb_user'@'localhost' IDENTIFIED VIA mysql_native_password USING PASSWORD 'password_anda_disini';

-- Berikan hak akses penuh ke database ecbtest
GRANT ALL PRIVILEGES ON ecbtest.* TO 'ecb_user'@'localhost';

-- Refresh hak akses
FLUSH PRIVILEGES;

-- Keluar
EXIT;
```

---

## 3. Persiapan Lingkungan Jalankan (Environment Setup)

### A. Instal Dependency Fyne (GUI Linux)
Aplikasi ini menggunakan Fyne untuk UI, sehingga butuh library grafis:
```bash
sudo apt install libgl1-mesa-dev xorg-dev libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev -y
```

### B. Instal Go (Golang)
Jika belum ada Go di Raspberry Pi:
```bash
# Cek arsitektur RPi (biasanya armv7l untuk RPi 3/4 32-bit atau aarch64 untuk 64-bit)
uname -m

# Contoh download Go 1.21.x (Sesuaikan dengan arsitektur Anda)
# Untuk ARMv7 (32-bit):
wget https://go.dev/dl/go1.21.6.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-armv6l.tar.gz

# Tambahkan ke Path
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

---

## 4. Menyiapkan Aplikasi

### A. Clone atau Copy Project
Copy folder project `V2go-ecbtest` ke Raspberry Pi.

### B. Konfigurasi .env
Copy file `.env.example` menjadi `.env` dan sesuaikan nilainya:
```bash
cp .env.example .env
nano .env
```
Pastikan bagian database sesuai:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=ecb_user
DB_PASSWORD=password_anda_disini
DB_DATABASE=ecbtest

# Untuk Raspberry Pi, set mode LIVE jika sudah terhubung ke GPIO
ECB_MODE=LIVE
```

### C. Jalankan Migrasi Database
Untuk membuat tabel-tabel yang diperlukan:
```bash
# Pastikan berada di root folder project
go run cmd/migrate/main.go up
```

---

## 5. Build dan Menjalankan Aplikasi

### A. Kompilasi (Build)
Untuk hasil terbaik di Raspberry Pi, gunakan target build yang sudah dioptimasi di Makefile:
```bash
make build-rpi-optimized
```
Hasil build akan berada di `bin/ecom-rpi`.

### B. Menjalankan Aplikasi
```bash
# Berikan izin eksekusi
chmod +x bin/ecom-rpi

# Jalankan aplikasi
./bin/ecom-rpi
```

---

## 6. Akses Localhost dan Auto-Start (Kiosk Mode)

Jika ingin aplikasi langsung terbuka otomatis saat Raspberry Pi menyala (seperti mesin ATM):

### A. Buat Launcher Script
Buat file `/home/pi/start_ecb.sh`:
```bash
#!/bin/bash
sleep 5
cd /home/pi/projects/new-goecbtest
export FYNE_RENDER=software
./bin/ecom-rpi
```
Jangan lupa berikan izin: `chmod +x /home/pi/start_ecb.sh`

### B. Konfigurasi Autostart
Buat file `~/.config/autostart/ecb_kiosk.desktop`:
```ini
[Desktop Entry]
Type=Application
Name=ECB Test Kiosk
Exec=/home/pi/start_ecb.sh
Terminal=false
```

### C. Matikan Screen Sleep
Gunakan `sudo raspi-config` -> **Display Options** -> **Screen Blanking** -> **No**.
Pastikan juga **Desktop Autologin** aktif di menu **System Options**.

---

## Troubleshooting Tips
- **Gagal Connect DB**: Cek status database dengan `sudo systemctl status mariadb`.
- **Error Auth Plugin**: Jika muncul error `plugin 'mysql_native_password' is not loaded`, pastikan saat CREATE USER menggunakan `IDENTIFIED VIA mysql_native_password`.
- **X11: The DISPLAY environment variable is missing**:
  Aplikasi ini membutuhkan GUI. Jika Anda menjalankan via SSH, gunakan perintah:
  ```bash
  export DISPLAY=:0
  ./bin/ecom-rpi
  ```
  Atau jalankan langsung melalui terminal di dalam desktop Raspberry Pi (bukan remote SSH).
- **Layar Putih/Blank**: Pastikan driver GPU Raspberry Pi aktif (cek `raspi-config` -> Advanced -> GL Driver).
