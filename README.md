# golox
An implementation of the [Lox language](https://craftinginterpreters.com/the-lox-language.html) written in go.

### Additional features
- **Multiline comments**
  
  ```
  /* This is
  a multiline 
  comment */
  ```
- **Anonymous functions**

  Anonymous functions for use cases like passing functions as arguments.
  ```
  fun thrice(fn) {
    for (var i = 1; i <= 3; i = i + 1) {
      fn(i);
    }
  }

  thrice(fun (a) {
    print a;
  });
  // "1".
  // "2".
  // "3".
  ```

## Usage
   Make sure you have [golang](https://go.dev/dl/) installed.  
  
   ```
   $ git clone https://github.com/AnshVM/golox.git
   $ cd ./golox
   $ go build
   $ ./golox
  ```
  This will fire up the lox REPl

  For running a lox file - 
  ```
  $ ./golox filepath.lox
  ```
