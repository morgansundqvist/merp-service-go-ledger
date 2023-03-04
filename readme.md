# Simplified ledger written in Go

## Limitations

- Voucher number solution does not work with horizontal scaling
- Only one voucher series per fiscal year
- All amounts are kept as int so user of api needs to transform int to precision decimal with two decimals
