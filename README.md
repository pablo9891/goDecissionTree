# Simple Go Decission Tree generator

Simple Decission Tree generator using ID3 algoritm.


* Applies ID3 algorithm to generate tree
* Handles csv files
* Handles different data types (numerical, boolean, character)
* Handles errors
* 3 data set samples
* Totally written in GO from scratch

## Examples

```bash
test@mac decissionTree % go run main.go "data_samples/buys_computer.csv" "buys_computer"
 +  age
  +  senior
   +  credit_rating
    +  excellent
     +  no
    +  fair
     +  yes
  +  middle_aged
   +  yes
  +  youth
   +  student
    +  yes
     +  yes
    +  no
     +  no
```

## Lessons learned

* Basic Golang syntax
* Pointers / strings library
* Go Modules and packages
* Interfaces and maps
* ID3 algorithm

## Future Work
* Handle undefined data types and empty values
* Improve code and delete useless functions
* Prunning the tree
* Configuration similar to SciKit learn algorithm

## References
* [Pandas Data Type Inference](https://rushter.com/blog/pandas-data-type-inference/)
* [Decission Tree in Machine Learning](https://kaumadiechamalka100.medium.com/decision-tree-in-machine-learning-c610ef087260#:~:text=ID3%20%28Iterative%20Dichotomiser%203%29%20Algorithm%20ID3%20%28Iterative%20Dichotomiser,and%20Information%20Gain%20to%20construct%20a%20decision%20tree.)
* [Decission Tree in SciKit Learn](https://scikit-learn.org/stable/modules/tree.html)
