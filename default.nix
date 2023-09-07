{buildGoModule}:
buildGoModule {
  pname = "goblin";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-TFSI51OEQrhbipjHligvPxDbxqpVWAwFyKoBjT+lql8=";

  ldflags = ["-s" "-w"];
}
