# Maintainer: TheRootDaemon <harsha.manjula.venkataramanan@gmail.com>

pkgname=tlgc
pkgver=1.0.0
pkgrel=1
pkgdesc="A tldr client written in Go"
arch=('x86_64' 'aarch64')
url="https://github.com/TheRootDaemon/tlgc"
license=('MIT')
depends=()
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::https://github.com/TheRootDaemon/tlgc/archive/refs/tags/v$pkgver.tar.gz")
sha256sums=('4da677a05956869b5bab353d459eb46c3b918c220f6f7efb4dac66e0127f2d82')

build() {
    cd "$pkgname-$pkgver"
    export CGO_ENABLED=0
    go build -trimpath \
        -ldflags "-s -w -X github.com/TheRootDaemon/tlgc/version.Version=$pkgver" \
        -o tldr ./cmd/tldr
}

package() {
    cd "$pkgname-$pkgver"
    install -Dm755 tldr "$pkgdir/usr/bin/tldr"
    install -Dm644 tldr.1 "$pkgdir/usr/share/man/man1/tldr.1"
    install -Dm644 completions/tldr.bash "$pkgdir/usr/share/bash-completion/completions/tldr"
    install -Dm644 completions/tldr.fish "$pkgdir/usr/share/fish/vendor_completions.d/tldr.fish"
    install -Dm644 completions/_tldr "$pkgdir/usr/share/zsh/site-functions/_tldr"
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
