for i in range(1, 8):
    # create temp var
    var = "tmp"
    for j in range(i):
        var += str(j)

    inp = var if i > 1 else "in[0]"
    out = "out" if i == 7 else var+str(i)

    # line up outputs
    spaces = (7 - i)*" "

    print(f"    Or(a={inp}, {spaces}b=in[{i}], out={out});");
