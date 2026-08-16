Add-Type -AssemblyName System.Drawing
$size = 512
$bmp = New-Object System.Drawing.Bitmap($size, $size)
$g = [System.Drawing.Graphics]::FromImage($bmp)
$g.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
$rect = New-Object System.Drawing.Rectangle(0, 0, $size, $size)
$brush = New-Object System.Drawing.Drawing2D.LinearGradientBrush($rect, [System.Drawing.Color]::FromArgb(10, 132, 255), [System.Drawing.Color]::FromArgb(94, 92, 230), 45)
$g.FillRectangle($brush, $rect)
$font = New-Object System.Drawing.Font("Segoe UI", 240, [System.Drawing.FontStyle]::Bold)
$format = New-Object System.Drawing.StringFormat
$format.Alignment = [System.Drawing.StringAlignment]::Center
$format.LineAlignment = [System.Drawing.StringAlignment]::Center
$white = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::White)
$g.DrawString("E", $font, $white, (New-Object System.Drawing.RectangleF(0, 0, $size, $size)), $format)
$outDir = Join-Path $PSScriptRoot "..\assets"
New-Item -ItemType Directory -Force -Path $outDir | Out-Null
$bmp.Save((Join-Path $outDir "icon.png"), [System.Drawing.Imaging.ImageFormat]::Png)
$g.Dispose(); $bmp.Dispose(); $brush.Dispose(); $font.Dispose(); $white.Dispose()
Write-Output "icon written to $(Join-Path $outDir 'icon.png')"
