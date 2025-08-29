{ pkgs, lib, config, inputs, ... }:

{
  # See full reference at https://devenv.sh/reference/options/

  # https://devenv.sh/packages/
  packages = [ pkgs.git ];

  # https://devenv.sh/languages/
  languages.go.enable = true;

  # https://devenv.sh/scripts/
  enterShell = ''
    git --version
    go version
  '';

  # https://devenv.sh/tests/
  # TODO: Go vet, Go staticcheck
  enterTest = ''
    echo "Running tests"
    git --version | grep --color=auto "${pkgs.git.version}"
  '';
}
