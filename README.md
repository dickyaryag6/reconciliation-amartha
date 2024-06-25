* Make sure you go to the folder of the repository
* Run the program by running this command `go build && ./amartha-test`
* In the code, i use port `8080`, you can change this accourding to your device, dont forget to also change the port in the code
* Use curl or postman to run test
  ```
  curl --location --request POST 'http://localhost:8000/reconciliation' \
  --form 'system_transactions=@"/Users/dickyarya/Documents/amartha-test/system_transactions.csv"' \
  --form 'bank_statements=@"/Users/dickyarya/Documents/amartha-test/bank_statements.csv"' 
  ```

* make sure to change value of `system_transaction` or `bank_statements` to the path of the file that you want to check
* make sure the format of the system transaction data is like the table below, or you can just look at the example file above :

trxID | amount | type | transactionTime 
--- | --- | --- | --- 
1 | "Rp8,500,000" | 2 | 01/01/2024 8:45:00
2 | "Rp7,000,000" | 2 | 01/01/2024 8:46:00
3 | "Rp6,000,000" | 2 | 01/01/2024 8:50:00
4 | "Rp8,500,000" | 1 | 03/01/2024 8:01:00
5 | "Rp7,000,000" | 2 | 04/01/2024 8:20:00
6 | "Rp6,000,000" | 2 | 04/01/2024 8:30:00
7 | "Rp8,500,000" | 1 | 08/01/2024 8:31:00
8 | "Rp7,000,000" | 2 | 12/01/2024 8:50:00
9 | "Rp6,000,000" | 2 | 12/01/2024 8:20:00
10 | "Rp2,000,000" | 2 | 15/02/2024 8:20:00

* make sure the format of the bank data is like the table below, or you can just look at the example file above :

unique_identifier | amount | date
--- | --- | --- 
BCA_12345 | "Rp1,500,000" | 01/01/2024
BRI_23463 | "Rp7,00,000" | 01/01/2024
BRI_23464 | "Rp6,000,000" | 01/01/2024
BRI_23465 | "-Rp8,500,000" | 03/01/2024
BRI_23466 | "Rp7,000,000" | 04/01/2024
BCA_12346 | "Rp6,000,000" | 04/01/2024
BCA_12347 | "-Rp8,500,000" | 08/01/2024
BCA_12348 | "-Rp7,000,000" | 12/01/2024
MANDIRI_12345 | "Rp6,000,000" | 12/01/2024
MANDIRI_12346 | "Rp1,500,000" | 13/01/2024
MANDIRI_12346 | "Rp2,500,000" | 13/01/2024

* amount or money have string data type in the code so the code could accept more than one currency, if there is any currency at all. Also to deal with formatting of money with two 0s behind the real number or with comma. If the amount or money value in csv is invalid or unknown, the code would return an error message that tells the user that the formatting is invalid

