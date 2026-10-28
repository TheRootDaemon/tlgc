{
  description = "A tldr client written in Go.";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = {
    self,
    nixpkgs,
  }: let
    # systems supported by this flake
    # each system gets its own package,
    # an runnable app and a development shell
    systems = [
      "x86_64-linux"
      "aarch64-linux"
      "aarch64-darwin"
    ];
    forAllSystems = nixpkgs.lib.genAttrs systems;

    # prefer a release tag when available,
    # if not a short git revision,
    # else "dev" when no information is available
    version = self.tag or (self.shortRev or "dev");
  in {
    # tools used when developing, testing, linting and releasing tlgc
    devShells = forAllSystems (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      default = pkgs.mkShell {
        packages = with pkgs; [
          go_1_27
          gofumpt
          golangci-lint
          gosec
          goreleaser
          scdoc
          actionlint
        ];
      };
    });

    # build applications separately for each supported system
    packages = forAllSystems (system: let
      pkgs = nixpkgs.legacyPackages.${system};
      tlgc = pkgs.buildGo127Module {
        src = self;
        vendorHash = "sha256-8Ir/RXt6Y4btJCos1SCCeCbWG5zBQdD8WpCMlT6Y5MY=";

        pname = "tlgc";
        inherit version;
        env = {CGO_ENABLED = 0;};

        doCheck = true;
        subPackages = ["cmd/tldr"];

        ldflags = [
          "-s"
          "-w"
          "-X"
          "github.com/TheRootDaemon/tlgc/version.Version=${version}"
        ];

        # install completions and the man page alongside the binary
        postInstall = ''
          install -Dm644 tldr.1 $out/share/man/man1/tldr.1
          install -Dm644 completions/tldr.bash $out/share/bash-completion/completions/tldr
          install -Dm644 completions/tldr.fish $out/share/fish/vendor_completions.d/tldr.fish
          install -Dm644 completions/_tldr $out/share/zsh/site-functions/_tldr
        '';

        meta = {
          description = "A tldr client written in Go";
          homepage = "https://github.com/TheRootDaemon/tlgc";
          license = pkgs.lib.licenses.mit;
          mainProgram = "tldr";
          platforms = pkgs.lib.platforms.unix;
        };
      };
    in {
      inherit tlgc;
      default = tlgc;
    });

    # define which executable `nix run` should launch for each system
    apps = forAllSystems (system: {
      default = {
        type = "app";
        program = "${self.packages.${system}.default}/bin/tldr";
      };
    });
  };
}
