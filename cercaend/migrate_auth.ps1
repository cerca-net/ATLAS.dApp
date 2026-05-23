$files = Get-ChildItem -Path "c:\Users\beatr\Desktop\ATLAS.dApp\cercaend\lib" -Filter *.dart -Recurse
foreach ($f in $files) {
    $content = Get-Content -Path $f.FullName -Raw
    if ($content -match "firebase_auth/auth_util\.dart") {
        $content = $content -replace "auth/firebase_auth/auth_util\.dart", "auth/auth_util.dart"
        $content = $content -replace "\.\./auth/firebase_auth/auth_util\.dart", "../auth/auth_util.dart"
        $content = $content -replace "firebase_auth/auth_util\.dart", "auth_util.dart"
        Set-Content -Path $f.FullName -Value $content -NoNewline
    }
}
