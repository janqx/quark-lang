def fib(n):
  return 1 if n < 3 else fib(n-1)+fib(n-2)

print(fib(35))