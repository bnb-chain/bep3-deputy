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

// Erc20SwapABI is the input ABI used to generate the binding from.
const Erc20SwapABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"ERC20ContractAddr\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"isSwapExist\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_randomNumberHash\",\"type\":\"bytes32\"},{\"name\":\"_swapSender\",\"type\":\"address\"},{\"name\":\"_bep2SenderAddr\",\"type\":\"bytes20\"}],\"name\":\"calSwapID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_randomNumber\",\"type\":\"bytes32\"}],\"name\":\"claim\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_randomNumberHash\",\"type\":\"bytes32\"},{\"name\":\"_timestamp\",\"type\":\"uint64\"},{\"name\":\"_heightSpan\",\"type\":\"uint256\"},{\"name\":\"_recipientAddr\",\"type\":\"address\"},{\"name\":\"_bep2SenderAddr\",\"type\":\"bytes20\"},{\"name\":\"_bep2RecipientAddr\",\"type\":\"bytes20\"},{\"name\":\"_outAmount\",\"type\":\"uint256\"},{\"name\":\"_bep2Amount\",\"type\":\"uint256\"}],\"name\":\"htlt\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"claimable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refundable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"queryOpenSwap\",\"outputs\":[{\"name\":\"_randomNumberHash\",\"type\":\"bytes32\"},{\"name\":\"_timestamp\",\"type\":\"uint64\"},{\"name\":\"_expireHeight\",\"type\":\"uint256\"},{\"name\":\"_outAmount\",\"type\":\"uint256\"},{\"name\":\"_sender\",\"type\":\"address\"},{\"name\":\"_recipient\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_erc20Contract\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_msgSender\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_recipientAddr\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_randomNumberHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_timestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"_bep2Addr\",\"type\":\"bytes20\"},{\"indexed\":false,\"name\":\"_expireHeight\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_outAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_bep2Amount\",\"type\":\"uint256\"}],\"name\":\"HTLT\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_msgSender\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_recipientAddr\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_randomNumberHash\",\"type\":\"bytes32\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_msgSender\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_recipientAddr\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_randomNumberHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_randomNumber\",\"type\":\"bytes32\"}],\"name\":\"Claimed\",\"type\":\"event\"}]"

// Erc20Swap is an auto generated Go binding around an Ethereum contract.
type Erc20Swap struct {
	Erc20SwapCaller     // Read-only binding to the contract
	Erc20SwapTransactor // Write-only binding to the contract
	Erc20SwapFilterer   // Log filterer for contract events
}

// Erc20SwapCaller is an auto generated read-only Go binding around an Ethereum contract.
type Erc20SwapCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20SwapTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Erc20SwapTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20SwapFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Erc20SwapFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20SwapSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Erc20SwapSession struct {
	Contract     *Erc20Swap        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc20SwapCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Erc20SwapCallerSession struct {
	Contract *Erc20SwapCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// Erc20SwapTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Erc20SwapTransactorSession struct {
	Contract     *Erc20SwapTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// Erc20SwapRaw is an auto generated low-level Go binding around an Ethereum contract.
type Erc20SwapRaw struct {
	Contract *Erc20Swap // Generic contract binding to access the raw methods on
}

// Erc20SwapCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Erc20SwapCallerRaw struct {
	Contract *Erc20SwapCaller // Generic read-only contract binding to access the raw methods on
}

// Erc20SwapTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Erc20SwapTransactorRaw struct {
	Contract *Erc20SwapTransactor // Generic write-only contract binding to access the raw methods on
}

// NewErc20Swap creates a new instance of Erc20Swap, bound to a specific deployed contract.
func NewErc20Swap(address common.Address, backend bind.ContractBackend) (*Erc20Swap, error) {
	contract, err := bindErc20Swap(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erc20Swap{Erc20SwapCaller: Erc20SwapCaller{contract: contract}, Erc20SwapTransactor: Erc20SwapTransactor{contract: contract}, Erc20SwapFilterer: Erc20SwapFilterer{contract: contract}}, nil
}

// NewErc20SwapCaller creates a new read-only instance of Erc20Swap, bound to a specific deployed contract.
func NewErc20SwapCaller(address common.Address, caller bind.ContractCaller) (*Erc20SwapCaller, error) {
	contract, err := bindErc20Swap(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20SwapCaller{contract: contract}, nil
}

// NewErc20SwapTransactor creates a new write-only instance of Erc20Swap, bound to a specific deployed contract.
func NewErc20SwapTransactor(address common.Address, transactor bind.ContractTransactor) (*Erc20SwapTransactor, error) {
	contract, err := bindErc20Swap(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20SwapTransactor{contract: contract}, nil
}

// NewErc20SwapFilterer creates a new log filterer instance of Erc20Swap, bound to a specific deployed contract.
func NewErc20SwapFilterer(address common.Address, filterer bind.ContractFilterer) (*Erc20SwapFilterer, error) {
	contract, err := bindErc20Swap(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Erc20SwapFilterer{contract: contract}, nil
}

// bindErc20Swap binds a generic wrapper to an already deployed contract.
func bindErc20Swap(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc20SwapABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20Swap *Erc20SwapRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Erc20Swap.Contract.Erc20SwapCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20Swap *Erc20SwapRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20Swap.Contract.Erc20SwapTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20Swap *Erc20SwapRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20Swap.Contract.Erc20SwapTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20Swap *Erc20SwapCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Erc20Swap.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20Swap *Erc20SwapTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20Swap.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20Swap *Erc20SwapTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20Swap.Contract.contract.Transact(opts, method, params...)
}

// ERC20ContractAddr is a free data retrieval call binding the contract method 0x49404437.
//
// Solidity: function ERC20ContractAddr() constant returns(address)
func (_Erc20Swap *Erc20SwapCaller) ERC20ContractAddr(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Erc20Swap.contract.Call(opts, out, "ERC20ContractAddr")
	return *ret0, err
}

// ERC20ContractAddr is a free data retrieval call binding the contract method 0x49404437.
//
// Solidity: function ERC20ContractAddr() constant returns(address)
func (_Erc20Swap *Erc20SwapSession) ERC20ContractAddr() (common.Address, error) {
	return _Erc20Swap.Contract.ERC20ContractAddr(&_Erc20Swap.CallOpts)
}

// ERC20ContractAddr is a free data retrieval call binding the contract method 0x49404437.
//
// Solidity: function ERC20ContractAddr() constant returns(address)
func (_Erc20Swap *Erc20SwapCallerSession) ERC20ContractAddr() (common.Address, error) {
	return _Erc20Swap.Contract.ERC20ContractAddr(&_Erc20Swap.CallOpts)
}

// CalSwapID is a free data retrieval call binding the contract method 0x7ef3e92e.
//
// Solidity: function calSwapID(bytes32 _randomNumberHash, address _swapSender, bytes20 _bep2SenderAddr) constant returns(bytes32)
func (_Erc20Swap *Erc20SwapCaller) CalSwapID(opts *bind.CallOpts, _randomNumberHash [32]byte, _swapSender common.Address, _bep2SenderAddr [20]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Erc20Swap.contract.Call(opts, out, "calSwapID", _randomNumberHash, _swapSender, _bep2SenderAddr)
	return *ret0, err
}

// CalSwapID is a free data retrieval call binding the contract method 0x7ef3e92e.
//
// Solidity: function calSwapID(bytes32 _randomNumberHash, address _swapSender, bytes20 _bep2SenderAddr) constant returns(bytes32)
func (_Erc20Swap *Erc20SwapSession) CalSwapID(_randomNumberHash [32]byte, _swapSender common.Address, _bep2SenderAddr [20]byte) ([32]byte, error) {
	return _Erc20Swap.Contract.CalSwapID(&_Erc20Swap.CallOpts, _randomNumberHash, _swapSender, _bep2SenderAddr)
}

// CalSwapID is a free data retrieval call binding the contract method 0x7ef3e92e.
//
// Solidity: function calSwapID(bytes32 _randomNumberHash, address _swapSender, bytes20 _bep2SenderAddr) constant returns(bytes32)
func (_Erc20Swap *Erc20SwapCallerSession) CalSwapID(_randomNumberHash [32]byte, _swapSender common.Address, _bep2SenderAddr [20]byte) ([32]byte, error) {
	return _Erc20Swap.Contract.CalSwapID(&_Erc20Swap.CallOpts, _randomNumberHash, _swapSender, _bep2SenderAddr)
}

// Claimable is a free data retrieval call binding the contract method 0x9b58e0a1.
//
// Solidity: function claimable(bytes32 _swapID) constant returns(bool)
func (_Erc20Swap *Erc20SwapCaller) Claimable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Erc20Swap.contract.Call(opts, out, "claimable", _swapID)
	return *ret0, err
}

// Claimable is a free data retrieval call binding the contract method 0x9b58e0a1.
//
// Solidity: function claimable(bytes32 _swapID) constant returns(bool)
func (_Erc20Swap *Erc20SwapSession) Claimable(_swapID [32]byte) (bool, error) {
	return _Erc20Swap.Contract.Claimable(&_Erc20Swap.CallOpts, _swapID)
}

// Claimable is a free data retrieval call binding the contract method 0x9b58e0a1.
//
// Solidity: function claimable(bytes32 _swapID) constant returns(bool)
func (_Erc20Swap *Erc20SwapCallerSession) Claimable(_swapID [32]byte) (bool, error) {
	return _Erc20Swap.Contract.Claimable(&_Erc20Swap.CallOpts, _swapID)
}

// IsSwapExist is a free data retrieval call binding the contract method 0x50f7a03b.
//
// Solidity: function isSwapExist(bytes32 _swapID) constant returns(bool)
func (_Erc20Swap *Erc20SwapCaller) IsSwapExist(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Erc20Swap.contract.Call(opts, out, "isSwapExist", _swapID)
	return *ret0, err
}

// IsSwapExist is a free data retrieval call binding the contract method 0x50f7a03b.
//
// Solidity: function isSwapExist(bytes32 _swapID) constant returns(bool)
func (_Erc20Swap *Erc20SwapSession) IsSwapExist(_swapID [32]byte) (bool, error) {
	return _Erc20Swap.Contract.IsSwapExist(&_Erc20Swap.CallOpts, _swapID)
}

// IsSwapExist is a free data retrieval call binding the contract method 0x50f7a03b.
//
// Solidity: function isSwapExist(bytes32 _swapID) constant returns(bool)
func (_Erc20Swap *Erc20SwapCallerSession) IsSwapExist(_swapID [32]byte) (bool, error) {
	return _Erc20Swap.Contract.IsSwapExist(&_Erc20Swap.CallOpts, _swapID)
}

// QueryOpenSwap is a free data retrieval call binding the contract method 0xb48017b1.
//
// Solidity: function queryOpenSwap(bytes32 _swapID) constant returns(bytes32 _randomNumberHash, uint64 _timestamp, uint256 _expireHeight, uint256 _outAmount, address _sender, address _recipient)
func (_Erc20Swap *Erc20SwapCaller) QueryOpenSwap(opts *bind.CallOpts, _swapID [32]byte) (struct {
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
	err := _Erc20Swap.contract.Call(opts, out, "queryOpenSwap", _swapID)
	return *ret, err
}

// QueryOpenSwap is a free data retrieval call binding the contract method 0xb48017b1.
//
// Solidity: function queryOpenSwap(bytes32 _swapID) constant returns(bytes32 _randomNumberHash, uint64 _timestamp, uint256 _expireHeight, uint256 _outAmount, address _sender, address _recipient)
func (_Erc20Swap *Erc20SwapSession) QueryOpenSwap(_swapID [32]byte) (struct {
	RandomNumberHash [32]byte
	Timestamp        uint64
	ExpireHeight     *big.Int
	OutAmount        *big.Int
	Sender           common.Address
	Recipient        common.Address
}, error) {
	return _Erc20Swap.Contract.QueryOpenSwap(&_Erc20Swap.CallOpts, _swapID)
}

// QueryOpenSwap is a free data retrieval call binding the contract method 0xb48017b1.
//
// Solidity: function queryOpenSwap(bytes32 _swapID) constant returns(bytes32 _randomNumberHash, uint64 _timestamp, uint256 _expireHeight, uint256 _outAmount, address _sender, address _recipient)
func (_Erc20Swap *Erc20SwapCallerSession) QueryOpenSwap(_swapID [32]byte) (struct {
	RandomNumberHash [32]byte
	Timestamp        uint64
	ExpireHeight     *big.Int
	OutAmount        *big.Int
	Sender           common.Address
	Recipient        common.Address
}, error) {
	return _Erc20Swap.Contract.QueryOpenSwap(&_Erc20Swap.CallOpts, _swapID)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(bytes32 _swapID) constant returns(bool)
func (_Erc20Swap *Erc20SwapCaller) Refundable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Erc20Swap.contract.Call(opts, out, "refundable", _swapID)
	return *ret0, err
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(bytes32 _swapID) constant returns(bool)
func (_Erc20Swap *Erc20SwapSession) Refundable(_swapID [32]byte) (bool, error) {
	return _Erc20Swap.Contract.Refundable(&_Erc20Swap.CallOpts, _swapID)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(bytes32 _swapID) constant returns(bool)
func (_Erc20Swap *Erc20SwapCallerSession) Refundable(_swapID [32]byte) (bool, error) {
	return _Erc20Swap.Contract.Refundable(&_Erc20Swap.CallOpts, _swapID)
}

// Claim is a paid mutator transaction binding the contract method 0x84cc9dfb.
//
// Solidity: function claim(bytes32 _swapID, bytes32 _randomNumber) returns(bool)
func (_Erc20Swap *Erc20SwapTransactor) Claim(opts *bind.TransactOpts, _swapID [32]byte, _randomNumber [32]byte) (*types.Transaction, error) {
	return _Erc20Swap.contract.Transact(opts, "claim", _swapID, _randomNumber)
}

// Claim is a paid mutator transaction binding the contract method 0x84cc9dfb.
//
// Solidity: function claim(bytes32 _swapID, bytes32 _randomNumber) returns(bool)
func (_Erc20Swap *Erc20SwapSession) Claim(_swapID [32]byte, _randomNumber [32]byte) (*types.Transaction, error) {
	return _Erc20Swap.Contract.Claim(&_Erc20Swap.TransactOpts, _swapID, _randomNumber)
}

// Claim is a paid mutator transaction binding the contract method 0x84cc9dfb.
//
// Solidity: function claim(bytes32 _swapID, bytes32 _randomNumber) returns(bool)
func (_Erc20Swap *Erc20SwapTransactorSession) Claim(_swapID [32]byte, _randomNumber [32]byte) (*types.Transaction, error) {
	return _Erc20Swap.Contract.Claim(&_Erc20Swap.TransactOpts, _swapID, _randomNumber)
}

// Htlt is a paid mutator transaction binding the contract method 0x91fda287.
//
// Solidity: function htlt(bytes32 _randomNumberHash, uint64 _timestamp, uint256 _heightSpan, address _recipientAddr, bytes20 _bep2SenderAddr, bytes20 _bep2RecipientAddr, uint256 _outAmount, uint256 _bep2Amount) returns(bool)
func (_Erc20Swap *Erc20SwapTransactor) Htlt(opts *bind.TransactOpts, _randomNumberHash [32]byte, _timestamp uint64, _heightSpan *big.Int, _recipientAddr common.Address, _bep2SenderAddr [20]byte, _bep2RecipientAddr [20]byte, _outAmount *big.Int, _bep2Amount *big.Int) (*types.Transaction, error) {
	return _Erc20Swap.contract.Transact(opts, "htlt", _randomNumberHash, _timestamp, _heightSpan, _recipientAddr, _bep2SenderAddr, _bep2RecipientAddr, _outAmount, _bep2Amount)
}

// Htlt is a paid mutator transaction binding the contract method 0x91fda287.
//
// Solidity: function htlt(bytes32 _randomNumberHash, uint64 _timestamp, uint256 _heightSpan, address _recipientAddr, bytes20 _bep2SenderAddr, bytes20 _bep2RecipientAddr, uint256 _outAmount, uint256 _bep2Amount) returns(bool)
func (_Erc20Swap *Erc20SwapSession) Htlt(_randomNumberHash [32]byte, _timestamp uint64, _heightSpan *big.Int, _recipientAddr common.Address, _bep2SenderAddr [20]byte, _bep2RecipientAddr [20]byte, _outAmount *big.Int, _bep2Amount *big.Int) (*types.Transaction, error) {
	return _Erc20Swap.Contract.Htlt(&_Erc20Swap.TransactOpts, _randomNumberHash, _timestamp, _heightSpan, _recipientAddr, _bep2SenderAddr, _bep2RecipientAddr, _outAmount, _bep2Amount)
}

// Htlt is a paid mutator transaction binding the contract method 0x91fda287.
//
// Solidity: function htlt(bytes32 _randomNumberHash, uint64 _timestamp, uint256 _heightSpan, address _recipientAddr, bytes20 _bep2SenderAddr, bytes20 _bep2RecipientAddr, uint256 _outAmount, uint256 _bep2Amount) returns(bool)
func (_Erc20Swap *Erc20SwapTransactorSession) Htlt(_randomNumberHash [32]byte, _timestamp uint64, _heightSpan *big.Int, _recipientAddr common.Address, _bep2SenderAddr [20]byte, _bep2RecipientAddr [20]byte, _outAmount *big.Int, _bep2Amount *big.Int) (*types.Transaction, error) {
	return _Erc20Swap.Contract.Htlt(&_Erc20Swap.TransactOpts, _randomNumberHash, _timestamp, _heightSpan, _recipientAddr, _bep2SenderAddr, _bep2RecipientAddr, _outAmount, _bep2Amount)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _swapID) returns(bool)
func (_Erc20Swap *Erc20SwapTransactor) Refund(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _Erc20Swap.contract.Transact(opts, "refund", _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _swapID) returns(bool)
func (_Erc20Swap *Erc20SwapSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _Erc20Swap.Contract.Refund(&_Erc20Swap.TransactOpts, _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _swapID) returns(bool)
func (_Erc20Swap *Erc20SwapTransactorSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _Erc20Swap.Contract.Refund(&_Erc20Swap.TransactOpts, _swapID)
}

// Erc20SwapClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the Erc20Swap contract.
type Erc20SwapClaimedIterator struct {
	Event *Erc20SwapClaimed // Event containing the contract specifics and raw log

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
func (it *Erc20SwapClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20SwapClaimed)
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
		it.Event = new(Erc20SwapClaimed)
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
func (it *Erc20SwapClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20SwapClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20SwapClaimed represents a Claimed event raised by the Erc20Swap contract.
type Erc20SwapClaimed struct {
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
func (_Erc20Swap *Erc20SwapFilterer) FilterClaimed(opts *bind.FilterOpts, _msgSender []common.Address, _recipientAddr []common.Address, _swapID [][32]byte) (*Erc20SwapClaimedIterator, error) {

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

	logs, sub, err := _Erc20Swap.contract.FilterLogs(opts, "Claimed", _msgSenderRule, _recipientAddrRule, _swapIDRule)
	if err != nil {
		return nil, err
	}
	return &Erc20SwapClaimedIterator{contract: _Erc20Swap.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0x9f46b1606087bdf4183ec7dfdbe68e4ab9129a6a37901c16a7b320ae11a96018.
//
// Solidity: event Claimed(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash, bytes32 _randomNumber)
func (_Erc20Swap *Erc20SwapFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *Erc20SwapClaimed, _msgSender []common.Address, _recipientAddr []common.Address, _swapID [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Erc20Swap.contract.WatchLogs(opts, "Claimed", _msgSenderRule, _recipientAddrRule, _swapIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20SwapClaimed)
				if err := _Erc20Swap.contract.UnpackLog(event, "Claimed", log); err != nil {
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
func (_Erc20Swap *Erc20SwapFilterer) ParseClaimed(log types.Log) (*Erc20SwapClaimed, error) {
	event := new(Erc20SwapClaimed)
	if err := _Erc20Swap.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// Erc20SwapHTLTIterator is returned from FilterHTLT and is used to iterate over the raw logs and unpacked data for HTLT events raised by the Erc20Swap contract.
type Erc20SwapHTLTIterator struct {
	Event *Erc20SwapHTLT // Event containing the contract specifics and raw log

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
func (it *Erc20SwapHTLTIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20SwapHTLT)
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
		it.Event = new(Erc20SwapHTLT)
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
func (it *Erc20SwapHTLTIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20SwapHTLTIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20SwapHTLT represents a HTLT event raised by the Erc20Swap contract.
type Erc20SwapHTLT struct {
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
func (_Erc20Swap *Erc20SwapFilterer) FilterHTLT(opts *bind.FilterOpts, _msgSender []common.Address, _recipientAddr []common.Address, _swapID [][32]byte) (*Erc20SwapHTLTIterator, error) {

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

	logs, sub, err := _Erc20Swap.contract.FilterLogs(opts, "HTLT", _msgSenderRule, _recipientAddrRule, _swapIDRule)
	if err != nil {
		return nil, err
	}
	return &Erc20SwapHTLTIterator{contract: _Erc20Swap.contract, event: "HTLT", logs: logs, sub: sub}, nil
}

// WatchHTLT is a free log subscription operation binding the contract event 0xb3e26d98380491276a8dce9d38fd1049e89070230ff5f36ebb55ead64500ade1.
//
// Solidity: event HTLT(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash, uint64 _timestamp, bytes20 _bep2Addr, uint256 _expireHeight, uint256 _outAmount, uint256 _bep2Amount)
func (_Erc20Swap *Erc20SwapFilterer) WatchHTLT(opts *bind.WatchOpts, sink chan<- *Erc20SwapHTLT, _msgSender []common.Address, _recipientAddr []common.Address, _swapID [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Erc20Swap.contract.WatchLogs(opts, "HTLT", _msgSenderRule, _recipientAddrRule, _swapIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20SwapHTLT)
				if err := _Erc20Swap.contract.UnpackLog(event, "HTLT", log); err != nil {
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
func (_Erc20Swap *Erc20SwapFilterer) ParseHTLT(log types.Log) (*Erc20SwapHTLT, error) {
	event := new(Erc20SwapHTLT)
	if err := _Erc20Swap.contract.UnpackLog(event, "HTLT", log); err != nil {
		return nil, err
	}
	return event, nil
}

// Erc20SwapRefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the Erc20Swap contract.
type Erc20SwapRefundedIterator struct {
	Event *Erc20SwapRefunded // Event containing the contract specifics and raw log

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
func (it *Erc20SwapRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20SwapRefunded)
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
		it.Event = new(Erc20SwapRefunded)
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
func (it *Erc20SwapRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20SwapRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20SwapRefunded represents a Refunded event raised by the Erc20Swap contract.
type Erc20SwapRefunded struct {
	MsgSender        common.Address
	RecipientAddr    common.Address
	SwapID           [32]byte
	RandomNumberHash [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0x04eb8ae268f23cfe2f9d72fa12367b104af16959f6a93530a4cc0f50688124f9.
//
// Solidity: event Refunded(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash)
func (_Erc20Swap *Erc20SwapFilterer) FilterRefunded(opts *bind.FilterOpts, _msgSender []common.Address, _recipientAddr []common.Address, _swapID [][32]byte) (*Erc20SwapRefundedIterator, error) {

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

	logs, sub, err := _Erc20Swap.contract.FilterLogs(opts, "Refunded", _msgSenderRule, _recipientAddrRule, _swapIDRule)
	if err != nil {
		return nil, err
	}
	return &Erc20SwapRefundedIterator{contract: _Erc20Swap.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0x04eb8ae268f23cfe2f9d72fa12367b104af16959f6a93530a4cc0f50688124f9.
//
// Solidity: event Refunded(address indexed _msgSender, address indexed _recipientAddr, bytes32 indexed _swapID, bytes32 _randomNumberHash)
func (_Erc20Swap *Erc20SwapFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *Erc20SwapRefunded, _msgSender []common.Address, _recipientAddr []common.Address, _swapID [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Erc20Swap.contract.WatchLogs(opts, "Refunded", _msgSenderRule, _recipientAddrRule, _swapIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20SwapRefunded)
				if err := _Erc20Swap.contract.UnpackLog(event, "Refunded", log); err != nil {
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
func (_Erc20Swap *Erc20SwapFilterer) ParseRefunded(log types.Log) (*Erc20SwapRefunded, error) {
	event := new(Erc20SwapRefunded)
	if err := _Erc20Swap.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	return event, nil
}
