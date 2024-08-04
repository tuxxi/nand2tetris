for i in range(16):
    if i < 10:
        spc = '   '
    else:
        spc = ' '
    print(f"And(a=a[{i}], b=b[{i}],{spc}out=out[{i}]);")
