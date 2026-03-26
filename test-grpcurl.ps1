# Konfigurasi
$HttpBaseUrl = "http://localhost:8081"
$GrpcAddress = "localhost:50051"
$ProtoFile = "proto/auth.proto" # Path yang benar

# Cek apakah grpcurl terinstall
if (-not (Get-Command grpcurl -ErrorAction SilentlyContinue)) {
    Write-Error "Tool 'grpcurl' tidak ditemukan. Silakan install: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    exit
}

Write-Host "--- 1. Testing Login via HTTP untuk mendapatkan Token ---" -ForegroundColor Cyan
$loginPayload = @{
    username = "admin_gudang"
    password = "password123"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$HttpBaseUrl/auth/login" -Method Post -Body $loginPayload -ContentType "application/json"
    $token = $response.data.token
    Write-Host "Login Berhasil! Token didapatkan." -ForegroundColor Green
    write-Host "response: $response"
} catch {
    Write-Error "Gagal login. Pastikan server berjalan dan user 'admin_gudang' sudah diregister via api_http.test"
    Write-Error $_
    exit
}

Write-Host "`n--- 2. Testing gRPC: GetUserProfile (ID=1) ---" -ForegroundColor Cyan
# Mengambil data user dengan ID 1
grpcurl -plaintext -proto $ProtoFile -d '{\"id\": 1}' $GrpcAddress auth.AuthService/GetUserProfile

Write-Host "`n--- 3. Testing gRPC: ValidateToken ---" -ForegroundColor Cyan
# Memvalidasi token yang barusan didapat dari login
# Kita harus menyusun JSON string dengan hati-hati untuk PowerShell
$jsonPayload = '{\"token\": \"' + $token + '\"}'

grpcurl -plaintext -proto $ProtoFile -d $jsonPayload $GrpcAddress auth.AuthService/ValidateToken

Write-Host "`nSelesai." -ForegroundColor Green