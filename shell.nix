{
  callPackage,
  gopls,
  go_1_21,
  revive,
}: let
  mainPkg = callPackage ./default.nix {};
in
  mainPkg.overrideAttrs (oa: {
    nativeBuildInputs =
      [
        gopls
        go_1_21
        revive
      ]
      ++ (oa.nativeBuildInputs or []);
  })
