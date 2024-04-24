<h2 align="center"> Tugas Besar 2 IF2211 Strategi Algoritma </h2>
<h1 align="center">  Wikirace by Dinasti</h1>

## Contributors
|   NIM    |                  Nama                  |
| :------: | :------------------------------------: |
| 13522009 |          Muhammad Yusuf Rafi           |
| 13522013 |        Denise Felicia Tiowanni         |


## Deskripsi Program
Program ini mengimplementasikan algoritma IDS dan BFS untuk menyelesaikan permainan WikiRace. Program diimplementasikan dalam bentuk sebuah website yang menerima masukan berupa jenis algoritma, judul artikel awal, dan judul artikel tujuan. Program kemudian memberikan keluaran berupa jumlah artikel yang diperiksa, jumlah artikel yang dilalui, rute penjelajahan artikel (dari artikel awal hingga artikel tujuan), dan waktu pencarian (dalam ms).

## Project Structure
```
│
├── doc
│   └── Dinasti.pdf
│
├── node_modules
│
├── src
│   ├── backend
│   │   ├── Api.go
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── log-bfs.txt
│   │   ├── log-ids.txt
│   └───└── dockerfile
│   ├── frontend
│   │   ├── wikirace
│   │   │   ├── node_modules
│   │   │   ├── public
│   │   │   ├── src
│   │   │   │   ├── components
│   │   │   │   │   ├── abous-us-page.js
│   │   │   │   │   ├── bfs-page.js
│   │   │   │   │   ├── button.js
│   │   │   │   │   ├── chopper-button.js
│   │   │   │   │   ├── header.js
│   │   │   │   │   ├── how-to-use-page.js
│   │   │   │   │   ├── ids-page.js
│   │   │   │   │   ├── luffy-button.js
│   │   │   │   │   └── styles.css.js
│   │   │   │   ├── App.js
│   │   │   │   └── index.js
│   │   │   ├── .gitignore
│   │   │   ├── dockerfile
│   │   │   ├── package-lock.json
│   └───└───└──package.json
│
├── docker-compose.yml
├── package-lock.json
├── package.json
├── tailwind.config.js
└── README.md

```

## Program Requirements

### Program requirements untuk backend (Go)
1.  Buka link <b>https://go.dev/doc/install</b> dan unduh package dengan menekan <code>download</code>. Selanjutnya, buka package file yang telah diinstal dan ikuti langkah penginstalan.

### Program requirements untuk frontend (React)
1. Install requirements website dengan command:
    ```
    npm install
    ```
2. Jalankan website dengan command:
    ```
    npm run start
    ```

### Program requirements jika ingin menjalankan dengan Docker
1. Pastikan telah menginstall Docker pada desktop. Untuk mendapatkan aplikasi dapat dengan menelusuri website <b>https://docs.docker.com/engine/install/</b>.


## How to Run
1. Clone repository ini dengan 
    ```
    git clone https://github.com/cupski/Tubes2_Dinasti.git
    ```
2. Buka folder repository pada terminal.
3. Pindah ke direktori *src* dengan `cd src`
4. Jika ingin menjalankan menggunakan Docker, buka aplikasi Docker pada deskrop. Lalu, masukkan command berikut pada terminal:
    ```
    docker-compose up --build
    ```
5. Buka <code>http://localhost:3000</code> pada browser dan website sudah dapat digunakan.
6. Apabila tidak ingin menjalankan menggunakan Docker, jalankan backend dan frontend secara terpisah.
7. Untuk menjalankan backend, Pindah ke direktori *backend* dengan `cd backend`.
8. Masukkan command berikut pada terminal:
    ```
    go run ./Api.go
    ```
9. Selanjutnya, untuk menjalankan frontend, buka terminal baru dan pindah ke direktori *wikirace* dengan `cd frontend/wikirace`.
10. Masukkan command berikut pada terminal:
    ```
    npm run start
    ```
11. Buka <code>http://localhost:3000</code> pada browser dan website sudah dapat digunakan.

## How to Use
1. Pilih metode penjelajahan yang diinginkan (BFS/IDS) dengan menekan tombol yang sesuai.
2. Masukkan start point yang diinginkan.
3. Masukkan end point yang diinginkan.
4. Klik tombol search.
5. Tunggu hasil keluar.
6. Refresh page jika ingin mencari lagi dan klik tombol back jika ingin kembali ke laman sebelumnya.
