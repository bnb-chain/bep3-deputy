// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ETHSwapABI is the input ABI used to generate the binding from.
const ETHSwapABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"isSwapExist\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_randomNumberHash\",\"type\":\"bytes32\"},{\"name\":\"_swapSender\",\"type\":\"address\"},{\"name\":\"_bep2SenderAddr\",\"type\":\"bytes20\"}],\"name\":\"calSwapID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_randomNumber\",\"type\":\"bytes32\"}],\"name\":\"claim\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"claimable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refundable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"queryOpenSwap\",\"outputs\":[{\"name\":\"_randomNumberHash\",\"type\":\"bytes32\"},{\"name\":\"_timestamp\",\"type\":\"uint64\"},{\"name\":\"_expireHeight\",\"type\":\"uint256\"},{\"name\":\"_outAmount\",\"type\":\"uint256\"},{\"name\":\"_sender\",\"type\":\"address\"},{\"name\":\"_recipient\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_randomNumberHash\",\"type\":\"bytes32\"},{\"name\":\"_timestamp\",\"type\":\"uint64\"},{\"name\":\"_heightSpan\",\"type\":\"uint256\"},{\"name\":\"_recipientAddr\",\"type\":\"address\"},{\"name\":\"_bep2SenderAddr\",\"type\":\"bytes20\"},{\"name\":\"_bep2RecipientAddr\",\"type\":\"bytes20\"},{\"name\":\"_bep2Amount\",\"type\":\"uint256\"}],\"name\":\"htlt\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_msgSender\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_recipientAddr\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_randomNumberHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_timestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"_bep2Addr\",\"type\":\"bytes20\"},{\"indexed\":false,\"name\":\"_expireHeight\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_outAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_bep2Amount\",\"type\":\"uint256\"}],\"name\":\"HTLT\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_msgSender\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_recipientAddr\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_randomNumberHash\",\"type\":\"bytes32\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_msgSender\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_recipientAddr\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_randomNumberHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_randomNumber\",\"type\":\"bytes32\"}],\"name\":\"Claimed\",\"type\":\"event\"}]"

// ETHSwap is an auto generated Go binding around an Ethereum contract.
type ETHSwap struct {
	ETHSwapCaller     // Read-only binding to the contract
	ETHSwapTransactor // Write-only binding to the contract
	ETHSwapFilterer   // Log filterer for contract events
}

// ETHSwapCaller is an auto generated read-only Go binding around an Ethereum contract.
type ETHSwapCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ETHSwapTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ETHSwapTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ETHSwapFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ETHSwapFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ETHSwapSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ETHSwapSession struct {
	Contract     *ETHSwap          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ETHSwapCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ETHSwapCallerSession struct {
	Contract *ETHSwapCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ETHSwapTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ETHSwapTransactorSession struct {
	Contract     *ETHSwapTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ETHSwapRaw is an auto generated low-level Go binding around an Ethereum contract.
type ETHSwapRaw struct {
	Contract *ETHSwap // Generic contract binding to access the raw methods on
}

// ETHSwapCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ETHSwapCallerRaw struct {
	Contract *ETHSwapCaller // Generic read-only contract binding to access the raw methods on
}

// ETHSwapTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ETHSwapTransactorRaw struct {
	Contract *ETHSwapTransactor // Generic write-only contract binding to access the raw methods on
}

// NewETHSwap creates a new instance of ETHSwap, bound to a specific deployed contract.
func NewETHSwap(address common.Address, backend bind.ContractBackend) (*ETHSwap, error) {
	contract, err := bindETHSwap(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ETHSwap{ETHSwapCaller: ETHSwapCaller{contract: contract}, ETHSwapTransactor: ETHSwapTransactor{contract: contract}, ETHSwapFilterer: ETHSwapFilterer{contract: contract}}, nil
}

// NewETHSwapCaller creates a new read-only instance of ETHSwap, bound to a specific deployed contract.
func NewETHSwapCaller(address common.Address, caller bind.ContractCaller) (*ETHSwapCaller, error) {
	contract, err := bindETHSwap(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ETHSwapCaller{contract: contract}, nil
}

// NewETHSwapTransactor creates a new write-only instance of ETHSwap, bound to a specific deployed contract.
func NewETHSwapTransactor(address common.Address, transactor bind.ContractTransactor) (*ETHSwapTransactor, error) {
	contract, err := bindETHSwap(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ETHSwapTransactor{contract: contract}, nil
}

// NewETHSwapFilterer creates a new log filterer instance of ETHSwap, bound to a specific deployed contract.
func NewETHSwapFilterer(address common.Address, filterer bind.ContractFilterer) (*ETHSwapFilterer, error) {
	contract, err := bindETHSwap(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ETHSwapFilterer{contract: contract}, nil
}

// bindETHSwap binds a generic wrapper to an already deployed contract.
func bindETHSwap(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ETHSwapABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ETHSwap *ETHSwapRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ETHSwap.Contract.ETHSwapCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ETHSwap *ETHSwapRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ETHSwap.Contract.ETHSwapTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ETHSwap *ETHSwapRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ETHSwap.Contract.ETHSwapTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ETHSwap *ETHSwapCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ETHSwap.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ETHSwap *ETHSwapTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ETHSwap.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ETHSwap *ETHSwapTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ETHSwap.Contract.contract.Transact(opts, method, params...)
}

// CalSwapID is a free data retrieval call binding the contract method 0x7ef3e92e.
//
// Solidity: function calSwapID(bytes32 _randomNumberHash, address _swapSender, bytes20 _bep2SenderAddr) constant returns(bytes32)
func (_ETHSwap *ETHSwapCaller) CalSwapID(opts *bind.CallOpts, _randomNumberHash [32]byte, _swapSender common.Address, _bep2SenderAddr [20]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _ETHSwap.contract.Call(opts, out, "calSwapID", _randomNumberHash, _swapSender, _bep2SenderAddr)
	return *ret0, err
}

// CalSwapID is a free data retrieval call binding the contract method 0x7ef3e92e.
//
// Solidity: function calSwapID(bytes32 _randomNumberHash, address _swapSender, bytes20 _bep2SenderAddr) constant returns(bytes32)
func (_ETHSwap *ETHSwapSession) CalSwapID(_randomNumberHash [32]byte, _swapSender common.Address, _bep2SenderAddr [20]byte) ([32]byte, error) {
	return _ETHSwap.Contract.CalSwapID(&_ETHSwap.CallOpts, _randomNumberHash, _swapSender, _bep2SenderAddr)
}

// CalSwapID is a free data retrieval call binding the contract method 0x7ef3e92e.
//
// Solidity: function calSwapID(bytes32 _randomNumberHash, address _swapSender, bytes20 _bep2SenderAddr) constant returns(bytes32)
func (_ETHSwap *ETHSwapCallerSession) CalSwapID(_randomNumberHash [32]byte, _swapSender common.Address, _bep2SenderAddr [20]byte) ([32]byte, error) {
	return _ETHSwap.Contract.CalSwapID(&_ETHSwap.CallOpts, _randomNumberHash, _swapSender, _bep2SenderAddr)
}

// Claimable is a free data retrieval call binding the contract method 0x9b58e0a1.
//
// Solidity: function claimable(bytes32 _swapID) constant returns(bool)
func (_ETHSwap *ETHSwapCaller) Claimable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ETHSwap.contract.Call(opts, out, "claimable", _swapID)
	return *ret0, err
}

// Claimable is a free data retrieval call binding the contract method 0x9b58e0a1.
//
// Solidity: function claimable(bytes32 _swapID) constant returns(bool)
func (_ETHSwap *ETHSwapSession) Claimable(_swapID [32]byte) (bool, error) {
	return _ETHSwap.Contract.Claimable(&_ETHSwap.CallOpts, _swapID)
}

// Claimable is a free data retrieval call binding the contract method 0x9b58e0a1.
//
// Solidity: function claimable(bytes32 _swapID) constant returns(bool)
func (_ETHSwap *ETHSwapCallerSession) Claimable(_swapID [32]byte) (bool, error) {
	return _ETHSwap.Contract.Claimable(&_ETHSwap.CallOpts, _swapID)
}

// IsSwapExist is a free data retrieval call binding the contract method 0x50f7a03b.
//
// Solidity: function isSwapExist(bytes32 _swapID) constant returns(bool)
func (_ETHSwap *ETHSwapCaller) IsSwapExist(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ETHSwap.contract.Call(opts, out, "isSwapExist", _swapID)
	return *ret0, err
}

// IsSwapExist is a free data retrieval call binding the contract method 0x50f7a03b.
//
// Solidity: function isSwapExist(bytes32 _swapID) constant returns(bool)
func (_ETHSwap *ETHSwapSession) IsSwapExist(_swapID [32]byte) (bool, error) {
	return _ETHSwap.Contract.IsSwapExist(&_ETHSwap.CallOpts, _swapID)
}

// IsSwapExist is a free data retrieval call binding the contract method 0x50f7a03b.
//
// Solidity: function isSwapExist(bytes32 _swapID) constant returns(bool)
func (_ETHSwap *ETHSwapCallerSession) IsSwapExist(_swapID [32]byte) (bool, error) {
	return _ETHSwap.Contract.IsSwapExist(&_ETHSwap.CallOpts, _swapID)
}

// QueryOpenSwap is a free data retrieval call binding the contract method 0xb48017b1.
//
// Solidity: function queryOpenSwap(bytes32 _swapID) constant returns(bytes32 _randomNumberHash, uint64 _timestamp, uint256 _expireHeight, uint256 _outAmount, address _sender, address _recipient)
func (_ETHSwap *ETHSwapCaller) QueryOpenSwap(opts *bind.CallOpts, _swapID [32]byte) (struct {
	RandomNumberHash [32]byte
	Timestamp        uint64
	ExpireHeight     *big.Int
	OutAmount        *big.Int
	Sender           common.Address
	Recipient        common.Address
}, error) {
	ret := new(struct {
		RandomNumberHash [32]byte
		Timestamp        uint64
		ExpireHeight     *big.Int
		OutAmount        *big.Int
		Sender           common.Address
		Recipient        common.Address
	})
	out := ret
	err := _ETHSwap.contract.Call(opts, out, "queryOpenSwap", _swapID)
	return *ret, err
}

// QueryOpenSwap is a free data retrieval call binding the contract method 0xb48017b1.
//
// Solidity: function queryOpenSwap(bytes32 _swapID) constant returns(bytes32 _randomNumberHash, uint64 _timestamp, uint256 _expireHeight, uint256 _outAmount, address _sender, address _recipient)
func (_ETHSwap *ETHSwapSession) QueryOpenSwap(_swapID [32]byte) (struct {
	RandomNumberHash [32]byte
	Timestamp        uint64
	ExpireHeight     *big.Int
	OutAmount        *big.Int
	Sender           common.Address
	Recipient        common.Address
}, error) {
	return _ETHSwap.Contract.QueryOpenSwap(&_ETHSwap.CallOpts, _swapID)
}

// QueryOpenSwap is a free data retrieval call binding the contract method 0xb48017b1.
//
// Solidity: function queryOpenSwap(bytes32 _swapID) constant returns(bytes32 _randomNumberHash, uint64 _timestamp, uint256 _expireHeight, uint256 _outAmount, address _sender, address _recipient)
func (_ETHSwap *ETHSwapCallerSession) QueryOpenSwap(_swapID [32]byte) (struct {
	RandomNumberHash [32]byte
	Timestamp        uint64
	ExpireHeight     *big.Int
	OutAmount        *big.Int
	Sender           common.Address
	Recipient        common.Address
}, error) {
	return _ETHSwap.Contract.QueryOpenSwap(&_ETHSwap.CallOpts, _swapID)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(bytes32 _swapID) constant returns(bool)
func (_ETHSwap *ETHSwapCaller) Refundable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ETHSwap.contract.Call(opts, out, "refundable", _swapID)
	return *ret0, err
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(bytes32 _swapID) constant returns(bool)
func (_ETHSwap *ETHSwapSession) Refundable(_swapID [32]byte) (bool, error) {
	return _ETHSwap.Contract.Refundable(&_ETHSwap.CallOpts, _swapID)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(bytes32 _swapID) constant returns(bool)
func (_ETHSwap *ETHSwapCallerSession) Refundable(_swapID [32]byte) (bool, error) {
	return _ETHSwap.Contract.Refundable(&_ETHSwap.CallOpts, _swapID)
}

// Claim is a paid mutator transaction binding the contract method 0x84cc9dfb.
//
// Solidity: function claim(bytes32 _swapID, bytes32 _randomNumber) returns(bool)
func (_ETHSwap *ETHSwapTransactor) Claim(opts *bind.TransactOpts, _swapID [32]byte, _randomNumber [32]byte) (*types.Transaction, error) {
	return _ETHSwap.contract.Transact(opts, "claim", _swapID, _randomNumber)
}

// Claim is a paid mutator transaction binding the contract method 0x84cc9dfb.
//
// Solidity: function claim(bytes32 _swapID, bytes32 _randomNumber) returns(bool)
func (_ETHSwap *ETHSwapSession) Claim(_swapID [32]byte, _randomNumber [32]byte) (*types.Transaction, error) {
	return _ETHSwap.Contract.Claim(&_ETHSwap.TransactOpts, _swapID, _randomNumber)
}

// Claim is a paid mutator transaction binding the contract method 0x84cc9dfb.
//
// Solidity: function claim(bytes32 _swapID, bytes32 _randomNumber) returns(bool)
func (_ETHSwap *ETHSwapTransactorSession) Claim(_swapID [32]byte, _randomNumber [32]byte) (*types.Transaction, error) {
	return _ETHSwap.Contract.Claim(&_ETHSwap.TransactOpts, _swapID, _randomNumber)
}

// Htlt is a paid mutator transaction binding the contract method 0xeba38663.
//
// Solidity: function htlt(bytes32 _randomNumberHash, uint64 _timestamp, uint256 _heightSpan, address _recipientAddr, bytes20 _bep2SenderAddr, bytes20 _bep2RecipientAddr, uint256 _bep2Amount) returns(bool)
func (_ETHSwap *ETHSwapTransactor) Htlt(opts *bind.TransactOpts, _randomNumberHash [32]byte, _timestamp uint64, _heightSpan *big.Int, _recipientAddr common.Address, _bep2SenderAddr [20]byte, _bep2RecipientAddr [20]byte, _bep2Amount *big.Int) (*types.Transaction, error) {
	return _ETHSwap.contract.Transact(opts, "htlt", _randomNumberHash, _timestamp, _heightSpan, _recipientAddr, _bep2SenderAddr, _bep2RecipientAddr, _bep2Amount)
}

// Htlt is a paid mutator transaction binding the contract method 0xeba38663.
//
// Solidity: function htlt(bytes32 _randomNumberHash, uint64 _timestamp, uint256 _heightSpan, address _recipientAddr, bytes20 _bep2SenderAddr, bytes20 _bep2RecipientAddr, uint256 _bep2Amount) returns(bool)
func (_ETHSwap *ETHSwapSession) Htlt(_randomNumberHash [32]byte, _timestamp uint64, _heightSpan *big.Int, _recipientAddr common.Address, _bep2SenderAddr [20]byte, _bep2RecipientAddr [20]byte, _bep2Amount *big.Int) (*types.Transaction, error) {
	return _ETHSwap.Contract.Htlt(&_ETHSwap.TransactOpts, _randomNumberHash, _timestamp, _heightSpan, _recipientAddr, _bep2SenderAddr, _bep2RecipientAddr, _bep2Amount)
}

// Htlt is a paid mutator transaction binding the contract method 0xeba38663.
//
// Solidity: function htlt(bytes32 _randomNumberHash, uint64 _timestamp, uint256 _heightSpan, address _recipientAddr, bytes20 _bep2SenderAddr, bytes20 _bep2RecipientAddr, uint256 _bep2Amount) returns(bool)
func (_ETHSwap *ETHSwapTransactorSession) Htlt(_randomNumberHash [32]byte, _timestamp uint64, _heightSpan *big.Int, _recipientAddr common.Address, _bep2SenderAddr [20]byte, _bep2RecipientAddr [20]byte, _bep2Amount *big.Int) (*types.Transaction, error) {
	return _ETHSwap.Contract.Htlt(&_ETHSwap.TransactOpts, _randomNumberHash, _timestamp, _heightSpan, _recipientAddr, _bep2SenderAddr, _bep2RecipientAddr, _bep2Amount)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _swapID) returns(bool)
func (_ETHSwap *ETHSwapTransactor) Refund(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _ETHSwap.contract.Transact(opts, "refund", _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _swapID) returns(bool)
func (_ETHSwap *ETHSwapSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _ETHSwap.Contract.Refund(&_ETHSwap.TransactOpts, _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _swapID) returns(bool)
func (_ETHSwap *ETHSwapTransactorSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _ETHSwap.Contract.Refund(&_ETHSwap.TransactOpts, _swapID)
}

// ETHSwapClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the ETHSwap contract.
type ETHSwapClaimedIterator struct {
	Event *ETHSwapClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ETHSwapClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ETHSwapClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ETHSwapClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ETHSwapClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ETHSwapClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ETHSwapClaimed represents a Claimed event raised by the ETHSwap contract.
type ETHSwapClaimed struct {
	MsgSender        common.Address
	RecipientAddr    common.Address
	SwapID           [32]byte
	RandomNumberHash [32]byte
	RandomNumber     [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0x9f46b1606087bdf4183ec7dfdbe68e4ab9129a6a37901c16a7b320ae11a96018.
//
// Solidity: event Claimed(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash, bytes32 _randomNumber)
func (_ETHSwap *ETHSwapFilterer) FilterClaimed(opts *bind.FilterOpts, _msgSender []common.Address, _recipientAddr []common.Address, _swapID [][32]byte) (*ETHSwapClaimedIterator, error) {

	var _msgSenderRule []interface{}
	for _, _msgSenderItem := range _msgSender {
		_msgSenderRule = append(_msgSenderRule, _msgSenderItem)
	}
	var _recipientAddrRule []interface{}
	for _, _recipientAddrItem := range _recipientAddr {
		_recipientAddrRule = append(_recipientAddrRule, _recipientAddrItem)
	}
	var _swapIDRule []interface{}
	for _, _swapIDItem := range _swapID {
		_swapIDRule = append(_swapIDRule, _swapIDItem)
	}

	logs, sub, err := _ETHSwap.contract.FilterLogs(opts, "Claimed", _msgSenderRule, _recipientAddrRule, _swapIDRule)
	if err != nil {
		return nil, err
	}
	return &ETHSwapClaimedIterator{contract: _ETHSwap.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0x9f46b1606087bdf4183ec7dfdbe68e4ab9129a6a37901c16a7b320ae11a96018.
//
// Solidity: event Claimed(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash, bytes32 _randomNumber)
func (_ETHSwap *ETHSwapFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *ETHSwapClaimed, _msgSender []common.Address, _recipientAddr []common.Address, _swapID [][32]byte) (event.Subscription, error) {

	var _msgSenderRule []interface{}
	for _, _msgSenderItem := range _msgSender {
		_msgSenderRule = append(_msgSenderRule, _msgSenderItem)
	}
	var _recipientAddrRule []interface{}
	for _, _recipientAddrItem := range _recipientAddr {
		_recipientAddrRule = append(_recipientAddrRule, _recipientAddrItem)
	}
	var _swapIDRule []interface{}
	for _, _swapIDItem := range _swapID {
		_swapIDRule = append(_swapIDRule, _swapIDItem)
	}

	logs, sub, err := _ETHSwap.contract.WatchLogs(opts, "Claimed", _msgSenderRule, _recipientAddrRule, _swapIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ETHSwapClaimed)
				if err := _ETHSwap.contract.UnpackLog(event, "Claimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseClaimed is a log parse operation binding the contract event 0x9f46b1606087bdf4183ec7dfdbe68e4ab9129a6a37901c16a7b320ae11a96018.
//
// Solidity: event Claimed(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash, bytes32 _randomNumber)
func (_ETHSwap *ETHSwapFilterer) ParseClaimed(log types.Log) (*ETHSwapClaimed, error) {
	event := new(ETHSwapClaimed)
	if err := _ETHSwap.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ETHSwapHTLTIterator is returned from FilterHTLT and is used to iterate over the raw logs and unpacked data for HTLT events raised by the ETHSwap contract.
type ETHSwapHTLTIterator struct {
	Event *ETHSwapHTLT // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ETHSwapHTLTIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ETHSwapHTLT)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ETHSwapHTLT)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ETHSwapHTLTIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ETHSwapHTLTIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ETHSwapHTLT represents a HTLT event raised by the ETHSwap contract.
type ETHSwapHTLT struct {
	MsgSender        common.Address
	RecipientAddr    common.Address
	SwapID           [32]byte
	RandomNumberHash [32]byte
	Timestamp        uint64
	Bep2Addr         [20]byte
	ExpireHeight     *big.Int
	OutAmount        *big.Int
	Bep2Amount       *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterHTLT is a free log retrieval operation binding the contract event 0xb3e26d98380491276a8dce9d38fd1049e89070230ff5f36ebb55ead64500ade1.
//
// Solidity: event HTLT(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash, uint64 _timestamp, bytes20 _bep2Addr, uint256 _expireHeight, uint256 _outAmount, uint256 _bep2Amount)
func (_ETHSwap *ETHSwapFilterer) FilterHTLT(opts *bind.FilterOpts, _msgSender []common.Address, _recipientAddr []common.Address, _swapID [][32]byte) (*ETHSwapHTLTIterator, error) {

	var _msgSenderRule []interface{}
	for _, _msgSenderItem := range _msgSender {
		_msgSenderRule = append(_msgSenderRule, _msgSenderItem)
	}
	var _recipientAddrRule []interface{}
	for _, _recipientAddrItem := range _recipientAddr {
		_recipientAddrRule = append(_recipientAddrRule, _recipientAddrItem)
	}
	var _swapIDRule []interface{}
	for _, _swapIDItem := range _swapID {
		_swapIDRule = append(_swapIDRule, _swapIDItem)
	}

	logs, sub, err := _ETHSwap.contract.FilterLogs(opts, "HTLT", _msgSenderRule, _recipientAddrRule, _swapIDRule)
	if err != nil {
		return nil, err
	}
	return &ETHSwapHTLTIterator{contract: _ETHSwap.contract, event: "HTLT", logs: logs, sub: sub}, nil
}

// WatchHTLT is a free log subscription operation binding the contract event 0xb3e26d98380491276a8dce9d38fd1049e89070230ff5f36ebb55ead64500ade1.
//
// Solidity: event HTLT(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash, uint64 _timestamp, bytes20 _bep2Addr, uint256 _expireHeight, uint256 _outAmount, uint256 _bep2Amount)
func (_ETHSwap *ETHSwapFilterer) WatchHTLT(opts *bind.WatchOpts, sink chan<- *ETHSwapHTLT, _msgSender []common.Address, _recipientAddr []common.Address, _swapID [][32]byte) (event.Subscription, error) {

	var _msgSenderRule []interface{}
	for _, _msgSenderItem := range _msgSender {
		_msgSenderRule = append(_msgSenderRule, _msgSenderItem)
	}
	var _recipientAddrRule []interface{}
	for _, _recipientAddrItem := range _recipientAddr {
		_recipientAddrRule = append(_recipientAddrRule, _recipientAddrItem)
	}
	var _swapIDRule []interface{}
	for _, _swapIDItem := range _swapID {
		_swapIDRule = append(_swapIDRule, _swapIDItem)
	}

	logs, sub, err := _ETHSwap.contract.WatchLogs(opts, "HTLT", _msgSenderRule, _recipientAddrRule, _swapIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ETHSwapHTLT)
				if err := _ETHSwap.contract.UnpackLog(event, "HTLT", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseHTLT is a log parse operation binding the contract event 0xb3e26d98380491276a8dce9d38fd1049e89070230ff5f36ebb55ead64500ade1.
//
// Solidity: event HTLT(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash, uint64 _timestamp, bytes20 _bep2Addr, uint256 _expireHeight, uint256 _outAmount, uint256 _bep2Amount)
func (_ETHSwap *ETHSwapFilterer) ParseHTLT(log types.Log) (*ETHSwapHTLT, error) {
	event := new(ETHSwapHTLT)
	if err := _ETHSwap.contract.UnpackLog(event, "HTLT", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ETHSwapRefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the ETHSwap contract.
type ETHSwapRefundedIterator struct {
	Event *ETHSwapRefunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ETHSwapRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ETHSwapRefunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ETHSwapRefunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ETHSwapRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ETHSwapRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ETHSwapRefunded represents a Refunded event raised by the ETHSwap contract.
type ETHSwapRefunded struct {
	MsgSender        common.Address
	RecipientAddr    common.Address
	SwapID           [32]byte
	RandomNumberHash [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0x04eb8ae268f23cfe2f9d72fa12367b104af16959f6a93530a4cc0f50688124f9.
//
// Solidity: event Refunded(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash)
func (_ETHSwap *ETHSwapFilterer) FilterRefunded(opts *bind.FilterOpts, _msgSender []common.Address, _recipientAddr []common.Address, _swapID [][32]byte) (*ETHSwapRefundedIterator, error) {

	var _msgSenderRule []interface{}
	for _, _msgSenderItem := range _msgSender {
		_msgSenderRule = append(_msgSenderRule, _msgSenderItem)
	}
	var _recipientAddrRule []interface{}
	for _, _recipientAddrItem := range _recipientAddr {
		_recipientAddrRule = append(_recipientAddrRule, _recipientAddrItem)
	}
	var _swapIDRule []interface{}
	for _, _swapIDItem := range _swapID {
		_swapIDRule = append(_swapIDRule, _swapIDItem)
	}

	logs, sub, err := _ETHSwap.contract.FilterLogs(opts, "Refunded", _msgSenderRule, _recipientAddrRule, _swapIDRule)
	if err != nil {
		return nil, err
	}
	return &ETHSwapRefundedIterator{contract: _ETHSwap.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0x04eb8ae268f23cfe2f9d72fa12367b104af16959f6a93530a4cc0f50688124f9.
//
// Solidity: event Refunded(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash)
func (_ETHSwap *ETHSwapFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *ETHSwapRefunded, _msgSender []common.Address, _recipientAddr []common.Address, _swapID [][32]byte) (event.Subscription, error) {

	var _msgSenderRule []interface{}
	for _, _msgSenderItem := range _msgSender {
		_msgSenderRule = append(_msgSenderRule, _msgSenderItem)
	}
	var _recipientAddrRule []interface{}
	for _, _recipientAddrItem := range _recipientAddr {
		_recipientAddrRule = append(_recipientAddrRule, _recipientAddrItem)
	}
	var _swapIDRule []interface{}
	for _, _swapIDItem := range _swapID {
		_swapIDRule = append(_swapIDRule, _swapIDItem)
	}

	logs, sub, err := _ETHSwap.contract.WatchLogs(opts, "Refunded", _msgSenderRule, _recipientAddrRule, _swapIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ETHSwapRefunded)
				if err := _ETHSwap.contract.UnpackLog(event, "Refunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRefunded is a log parse operation binding the contract event 0x04eb8ae268f23cfe2f9d72fa12367b104af16959f6a93530a4cc0f50688124f9.
//
// Solidity: event Refunded(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash)
func (_ETHSwap *ETHSwapFilterer) ParseRefunded(log types.Log) (*ETHSwapRefunded, error) {
	event := new(ETHSwapRefunded)
	if err := _ETHSwap.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	return event, nil
}
