package account

import "sync"

// Define the Account type here.
type Account struct {
	isOpen  bool
	balance int64
	mx      sync.Mutex
}

func Open(amount int64) *Account {
	if amount < 0 {
		return nil
	}
	a := Account{
		balance: amount,
		isOpen:  true,
	}
	return &a
}

func (a *Account) Balance() (int64, bool) {
	a.mx.Lock()
	defer a.mx.Unlock()
	if !a.isOpen {
		return 0, false
	}
	return a.balance, true
}

func (a *Account) Deposit(amount int64) (int64, bool) {
	a.mx.Lock()
	defer a.mx.Unlock()

	if !a.isOpen {
		return 0, false
	}

	if a.balance+amount >= 0 {
		a.balance += amount
		return a.balance, true
	}
	return 0, false
}

func (a *Account) Close() (int64, bool) {
	a.mx.Lock()
	defer a.mx.Unlock()

	if !a.isOpen {
		return 0, false
	}

	b := a.balance
	a.balance = 0
	a.isOpen = false
	return b, true
}
