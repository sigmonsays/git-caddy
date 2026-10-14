{ pkgs ? import <nixpkgs> {} }:

with pkgs;

mkShell {
  buildInputs = [

    git
    gomod2nix

    # Go development
    go_1_27
    gopls
    godef
    gotools

  ];
}
