{
  inputs = {
    nixpkgs = {
      url = "github:NixOS/nixpkgs/nixos-unstable";
    };
    flake-utils = {
      url = "github:numtide/flake-utils";
    };
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };

        # Define the exact native compilation and graphics libraries Fyne expects
        buildDeps = with pkgs; [
          pkg-config
        ];

        runtimeDeps = with pkgs; [
          libX11
          libXcursor
          libXrandr
          libXinerama
          libXi
          libXext
          libXfixes
          libGL
          libxkbcommon
          wayland
          libXxf86vm
        ];
      in
      {
        packages = {
          default = pkgs.buildGoModule {
            pname = "alwaysFocusedInputBox";
            version = "2";
            src = ./.;
            vendorHash = "sha256-Yl4RTwNsmVuXWerWgRDwq4fFbrTxFzd7sUv6N9aq8PM=";

            nativeBuildInputs = buildDeps;
            buildInputs = runtimeDeps;
          };
        };

        devShells = {
          default = pkgs.mkShell {
            nativeBuildInputs = buildDeps;
            buildInputs =
              with pkgs;
              [
                go
                gopls
              ]
              ++ runtimeDeps;

            # Crucial for NixOS: Tells Cgo exactly where to find your live graphics drivers
            shellHook = ''
              export LD_LIBRARY_PATH=${pkgs.lib.makeLibraryPath runtimeDeps}:$LD_LIBRARY_PATH
            '';
          };
        };
      }
    );
}
