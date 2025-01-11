let x = 5;
let y = 10;
let z = 20;

let arr = [x, y, z];
for (let i = 0; i < len(arr); i++) {
    puts("arr[", i, "] is: ", arr[i]);
}
puts("")

let adder = fn(x) {
    return x + 1;
}

i = 0
while (i < len(arr)) {
    puts("arr[", i, "] + 1 is: ", adder(arr[i]));
    i++;
}
puts("")

puts("The sum of elements in arr is: ", sum(arr));
puts("")
