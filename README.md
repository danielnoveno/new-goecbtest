# Go ECB Test

Aplikasi desktop berbasis Fyne untuk pengujian Electronic Control Board (ECB) di lini produksi. UI menampilkan status line, menerima input scan (SN/SPC/kompresor), serta mengelola login, tema, dan bahasa. Seluruh backend berjalan di dalam proses desktop: `services/gpio` berbicara langsung ke simulator atau Raspberry Pi melalui driver GPIO, sementara scheduler menyinkronkan master data dan mengirim hasil scan ke database eksternal tanpa membutuhkan server HTTP terpisah.

## Fitur Utama
- UI desktop: sidebar dinamis, tema ganti langsung, pilih bahasa (id/en), form login/register, halaman profil dan perawatan.
- Mode line fleksibel: `sn-only-*`, `refrig-*`, `refrig-po-*` dengan dukungan single/double line dan line selection.
- Kontrol GPIO lokal: panel Maintenance memicu `services/gpio` (Initialize, Start, Reset, Line Select) dan polling status langsung ke driver tanpa HTTP.
- Scheduler: sync master dari SIMO, tarik PO, push `ecbdatas` ke SIMO & bserv, bersihkan mutex jadwal.
- Migrasi & seed: skrip SQL di `cmd/migrate/migrations` dan dummy data di `db/demo_data.sql`.
- Build lintas platform: target default (desktop), RPi/ARM, dan paket executable Fyne.

## Prasyarat
- Go 1.24+.
- MySQL/MariaDB untuk database lokal dan (opsional) koneksi SIMO/BSERV.
- (Opsional) `make`, `golang-migrate`, dan toolchain Fyne untuk packaging GUI (`fyne package`).

## Setup Cepat
1) Salin dan lengkapi konfigurasi  
   - Gunakan `.env` yang tersedia atau duplikasi ke `.env.local` sesuai kebutuhan.  
   - Isi kredensial database lokal, serta host/user/password untuk SIMO & bserv jika ingin sync/POST data.

2) Unduh dependensi  
   ```bash
   go mod download
   ```

3) Siapkan schema  
   - Cara cepat: jalankan aplikasi sekali supaya gorp membuat tabel dasar.  
     ```bash
     go run ./cmd
     ```  
   - Atau gunakan migrasi:  
     ```bash
     go run cmd/migrate/main.go up
     # make migrate-up  (opsi Makefile)
     ```

4) (Opsional) Seed dummy data untuk demo  
   ```bash
   mysql -u "$DB_USERNAME" -p"$DB_PASSWORD" "$DB_DATABASE" < db/demo_data.sql
   ```
   Kredensial default dummy: `admin@example.com / password`.

5) Jalankan aplikasi  
   ```bash
   go run ./cmd           # mode dev cepat
   # atau
   make run               # build + eksekusi bin/ecom
   ```
   Aplikasi hanya berjalan sebagai desktop; perintah GPIO dan simulasi ditrigger dari menu Maintenance. Tidak ada API HTTP yang harus diakses oleh pengguna.

## Perintah Build & Deploy
- `make build`             : build desktop binary ke `bin/ecom`.
- `make build-rpi`         : build target linux/arm (sesuaikan `RPI_GOOS`, `RPI_GOARCH`, `RPI_GOARM`).
- `make fyne-build`        : paket executable Windows (`bin/ECB Test.exe`).
- `make fyne-build-rpi`    : paket tar.xz untuk Linux ARM.

## Konfigurasi (.env)
Nilai bawaan contoh ada di `.env`. Semua boleh diset lewat env var sistem.

- Database lokal:
  - `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, `DB_DATABASE`
- Database SIMO (sinkronisasi master & kirim data):
  - `DBSIMOPRD_HOST`, `DBSIMOPRD_PORT`, `DBSIMOPRD_DATABASE`, `DBSIMOPRD_USERNAME`, `DBSIMOPRD_PASSWORD`
- Database bserv (kirim data ecb):
  - `DBBSERV_HOST`, `DBBSERV_PORT`, `DBBSERV_DATABASE`, `DBBSERV_USERNAME`, `DBBSERV_PASSWORD`
- Auth/JWT:
  - `JWT_SECRET`, `JWT_EXPIRATION_IN_SECONDS`
- Info aplikasi & tampilan:
  - `APP_ENV`, `APP_DEBUG`, `APP_KEY`
  - `APP_TITLE`, `APP_VERSION`, `APP_DESCRIPTION`, `APP_ICON`
  - `APP_THEME_DEFAULT`   : nama tema dari tabel `themes` (fallback ke tema default jika kosong)
- Konfigurasi lini ECB:
  - `ECB_LOCATION`        : teks lokasi di header.
  - `ECB_LINE_TYPE`       : `sn-only-single`, `sn-only-double`, `refrig-single`, `refrig-double`, `refrig-po-single`, `refrig-po-double`.
  - `ECB_LINEIDS`         : daftar line, pisahkan koma (contoh `L1, L2` atau `REF A, REF B`).
  - `ECB_WORKCENTERS`     : kode workcenter, pisahkan koma.
  - `ECB_TACKTIME`        : takt time (detik).
  - `ECB_STATE_DEFAULT`   : pola state awal/SSE (contoh `0.0.0.1`).
  - `ECB_MODE`            : `LIVE`, `simulateAll`, `simulateHW`, `simulateDB`.
  - `EcbMode` mempengaruhi driver GPIO: `LIVE` membaca pin Raspberry Pi, `simulate*` menggunakan loop internal.

## Mode Simulasi vs Hardware
- `simulateAll`   : tanpa GPIO, status dipol oleh panel Maintenance mengikuti `ECB_STATE_DEFAULT` (aman untuk laptop/demo).
- `simulateHW`    : pin tidak dipakai, tapi handler tetap mengemulasi state untuk menjaga flow UI.
- `simulateDB`    : state dibaca dari tabel `ecbstates` untuk memutar log lama.
- `LIVE`          : jalankan di Raspberry Pi dengan wiring sesuai driver di `services/gpio`; Maintenance panel akan langsung mengontrol pin hardware.

## Kontrol GPIO Lokal
Perintah GPIO dijalankan langsung dari menu Maintenance. Panel tersebut:
- memanggil `gpio.InitializeControl`, `gpio.StartTest`, `gpio.ResetTest`, dan `gpio.LineToggle/LineSet` tanpa keluar dari proses desktop.
- Memonitor status pin lewat polling (`gpio.ReadLocalEcbState`) setiap 250 ms, sehingga UI menampilkan PASS/FAIL/UNDERTEST/line aktif secara real-time tanpa SSE atau HTTP streaming.
Aktivitas utama:
1. **RE-INIT** – konfigurasi ulang pin (semua output/input diset ulang dan lampu berkedip untuk tanda siap).
2. **START/RESET** – jalankan siklus pengujian yang sama seperti driver langsung: WRITE ke pin 23/24/28/29 dan baca pin 2/21/22/25.
3. **Line Select** – toggle atau langsung set line aktif, nilai `lineSelect` dibaca dari pin 25.
4. **Mode sim/Live** – pengaturan `ECB_MODE` di UI memengaruhi apakah driver berinteraksi dengan GPIO nyata atau hanya simulasi.

## Contoh Input Scan Saat Migrasi Awal

Setelah menjalankan `go run cmd/migrate/main.go up`, semua tabel utama masih kosong sampai seed tambahan dijalankan. Karena ini fase migrasi awal, asumsi kita cuma ada tabel `themes`/`navigations` hasil seed bawaan; data master lain (masterfgs, compressors, comprefgs, ecbpos) harus dimasukkan dulu agar validasi di `views/ecb/flows.go` tidak menolak input.

### Persiapan pertama

- Pastikan `.env` mengarahkan `ECB_LINE_TYPE` dan `ECB_LINE_IDS` sesuai mode yang ingin diuji, serta `ECB_WORKCENTERS` mencakup `WC-REF` agar `LineState` menulis ke workcenter yang sama dengan contoh data. Contoh kombinasi:
  - `ECB_LINE_TYPE=sn-only-single` → `ECB_LINE_IDS=REF A`, `ECB_WORKCENTERS=WC-REF`.
  - `ECB_LINE_TYPE=refrig-double` atau `refrig-po-double` → `ECB_LINE_IDS=REF A, REF B` dan `ECB_WORKCENTERS=WC-REF,WC-REF`.
  - Jalankan UI dengan `ECB_MODE=simulateAll` untuk bisa menguji tanpa perangkat keras.
- Jalankan `db/demo_data.sql` setelah migrasi untuk menambahkan semua master/compressed/P.O. yang dipakai oleh contoh—skrip ini juga mencantumkan `ecbdatas` final agar Anda bisa melihat hasil yang valid.

### Validasi yang terjadi per mode

1. **SN Only Single (`sn-only-single`)**
   - `makeSnOnlySerialValidator` mengecek 4 digit awal S/N ada di `masterfgs.kdbar`, `Lotinv = IDN0` memaksa panjang S/N harus 12 karakter, dan `ecbdatas.sn` tidak boleh pernah muncul sebelumnya.
   - Skrip contoh menyertakan satu row `ecbdatas` untuk menunjukkan format entri (ctgr `SNONLY`, `spc`/kompresor kosong).

2. **Refrig Single (`refrig-single`)**
   - `validateSpc` mengharuskan SPC 11 karakter, belum pernah tercatat lokal maupun pada server remote (`remote.exists`), agar proses tidak diteruskan.
   - `validateSerial` memeriksa `masterfgs`, memastikan panjang S/N untuk `Lotinv=IDN0`, dan tidak ada duplikat `ecbdatas.sn`.
   - `validateCompressorType` mencocokkan scanner dengan record `compressors` lalu memastikan prefix S/N ada di `comprefgs`.
   - `validateCompressorCode` mengecek bahwa kode kompresor tidak sekadar ulang-ulang string tipe jika `force_scan=1`.

3. **Refrig Double (`refrig-double`)**
   - Semua validasi sama seperti Refrig Single, tetapi `ECB_LINE_IDS` harus berisi dua line aktif sehingga dua kartu bisa dibaca (contoh `REF A, REF B` di `.env`). Skrip juga menyertakan dua baris `ecbdatas` dengan line yang berbeda untuk menggambarkan batch double-line.

4. **Refrig PO (`refrig-po-single` / `refrig-po-double`)**
   - `validateCompressorType` akan mengambil PO yang cocok dari tabel `ecbpos` (S/N + tipe kompresor harus sesuai) dan menyimpan referensi PO tersebut.
   - Saat `saveData` dipanggil, row baru di `ecbdatas` menyertakan kolom `po`, dan status `ecbpos` otomatis diperbarui menjadi `scanned`.

### Menjalankan contoh data

- File `db/demo_data.sql` berisi semua perintah SQL di atas plus master/compressor/PO yang diperlukan. Jalankan sekali setelah migrasi, lalu jalankan UI untuk mode yang sesuai.
- Jika ingin menambah prefix lain atau workcenter berbeda, cukup tambahkan baris di `masterfgs`, `comprefgs`, dan `ecbdatas` di skrip tersebut.
- Contoh data tidak menghapus tabel lain (themes, navigations, dsb.), jadi Anda tidak perlu menyentuh seed bawaan selain memastikan `.env` mengarah ke line/workcenter yang sama.

## Scheduler & Integrasi Data
Saat aplikasi jalan, scheduler (`task.Start`) otomatis:
- Membersihkan file mutex lama (`CleanMutex`).
- Menarik master dari SIMO (`ecbstations`, `masterfgs`, `mastersfgs`, `compressors`, `comprefgs`).
- Menarik PO (`SyncEcbPo`) dan mem-push `ecbdatas` ke SIMO serta bserv setiap beberapa menit.
- Menonaktifkan job eksternal jika host SIMO/bserv tidak bisa diping untuk mencegah macet UI.

## Lokal & Tema
- Terjemahan ada di `lang/id.json` dan `lang/en.json`. Dropdown bahasa di header mengubah teks runtime.
- Tema diambil dari tabel `themes` atau `APP_THEME_DEFAULT`; fallback ke `DefaultPalette` bila kosong.
- Logo bisa diganti lewat `APP_ICON` (path relatif/absolut; default `assets/logo-nb.webp`).

## Troubleshooting Ringkas
- Error lock: pastikan tidak ada instance lain; file lock di `%TEMP%/go-ecb.lock`.
- Build Fyne di Windows kadang butuh toolchain C (mingw) dan dependencies Fyne; jalankan `go env -w CGO_ENABLED=1` bila diperlukan.
- Jika scheduler gagal (log `[scheduler] disable SIMO/BServ jobs`), cek koneksi/credential SIMO & bserv.

## Struktur Penting
- `cmd/main.go`       : entrypoint desktop (Fyne), inisialisasi GPIO/scheduler, dan bangun UI utama.
- `cmd/migrate`       : runner migrasi golang-migrate.
- `services/gpio`     : driver/simulasi GPIO plus helpers lokal (start/reset, polling status, line select) untuk Maintenance panel.
- `views/*`           : komponen UI Fyne (auth, home, ecb, theme).
- `task/*`            : scheduler sync & push data.
- `db/demo_data.sql`  : seed dummy.
- `configs/environment_loader.go`    : loader env var `.env`.
