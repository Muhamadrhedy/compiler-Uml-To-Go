# Tahap 3 — Semantic Analysis

Semantic analyzer menerima AST dari parser dan memeriksa aturan yang tidak cukup diperiksa oleh grammar saja.

## Aturan MVP

1. Nama class tidak boleh dideklarasikan lebih dari sekali.
2. Nama atribut dan method harus unik dalam satu class. Overloading method belum didukung.
3. Nama parameter harus unik dalam satu method.
4. Tipe yang didukung adalah `int`, `float`, `string`, `bool`, dan `void`.
5. `void` hanya valid sebagai tipe kembalian method, bukan tipe atribut atau parameter.
6. Semua diagnostic menyertakan nomor baris jika tersedia.

## API

```go
errors := semantic.Check(program)
if len(errors) > 0 {
    // tampilkan diagnostic dan hentikan code generation
}
```

Pemeriksaan ini mengembalikan seluruh diagnostic yang ditemukan agar pengguna dapat memperbaiki beberapa masalah dalam satu proses. Parser tetap bertanggung jawab terhadap sintaks, sedangkan semantic analyzer memeriksa konsistensi AST.

## Jalankan pengujian

```bash
go test ./...
go vet ./...
```
