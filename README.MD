# Pepi Coin
This is my Golang implementation of a blockchain.  
  
## Big Idea
Compile this to C then to webassembly.  
Distribute and let everyone mine from their browser.  

## Bigger Idea 
Distubuted computing from a web of browser applications.  

## implementation

### transaction
simple implementation.  
timestamp - when the transaction was made  
from - address string  
to - address string  
amount - how much is being sent  
signature - signed by sender using from  

### block
prev - hash of previous block  
trasactions - up to 8 transactions in each block  
nonce - used for hashing  
hash - this block's hash  

### proof (of work)
keeps hashing until a certain number of 0 appear at the start of the hash  

### blockchain
stores each chain inside badger db  
initiates db with genesis block  

### wallet
geneates key pair using random words  

### network 
TBC  