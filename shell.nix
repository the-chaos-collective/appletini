{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell rec {
    nativeBuildInputs = [
        pkg-config
    ];
    buildInputs = [
        libGL xorg.libX11.dev xorg.libXcursor xorg.libXi xorg.libXinerama xorg.libXrandr xorg.libXxf86vm
    ];
}
