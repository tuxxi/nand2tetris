for i in range(1, 16):
    s = "    " if i < 10 else " "
    print(f"FullAdder(a=a[{i}], b=b[{i}], c=c{i},{s}sum=out[{i}], carry=c{i+1});");
