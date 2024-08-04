for i in range(16):
    if i < 10:
        spc = '   '
    else:
        spc = ' '
    print(f"Mux(a=a[{i}], b=b[{i}],sel=sel,{spc}out=out[{i}]);")

