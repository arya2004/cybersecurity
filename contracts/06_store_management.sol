// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract StoreManagement {
    
    // Enum to represent the status of the store
    enum StoreStatus { Open, Closed }
    StoreStatus public storeStatus;
    
    // Struct to represent a Product
    struct Product {
        string name;
        uint price; // Price in Wei
        uint stock; // Number of items available
    }

    // Array to store the list of products
    Product[] public products;

    // Variable to store total products sold
    uint public totalProductsSold;

    // Address of the store owner (admin)
    address public owner;

    // Constructor to initialize the store (opens the store and sets the owner)
    constructor() {
        owner = msg.sender;
        storeStatus = StoreStatus.Open;
        totalProductsSold = 0;
    }

    // Modifier to check if the caller is the owner
    modifier onlyOwner() {
        require(msg.sender == owner, "Only the owner can perform this action");
        _;
    }

    // Modifier to check if the store is open
    modifier storeOpen() {
        require(storeStatus == StoreStatus.Open, "Store is currently closed");
        _;
    }

    // Function to add a new product to the store (only owner can add products)
    function addProduct(string memory _name, uint _price, uint _stock) public onlyOwner storeOpen {
        products.push(Product({
            name: _name,
            price: _price,
            stock: _stock
        }));
    }

    // Function to buy a product by providing its index in the array
    function buyProduct(uint _productIndex) public payable storeOpen {
        require(_productIndex < products.length, "Product does not exist");
        Product storage product = products[_productIndex];

        // Decision Making: Check if the buyer sent enough ether and if the product is in stock
        require(msg.value >= product.price, "Not enough ether sent");
        require(product.stock > 0, "Product out of stock");

        // Decrease the stock of the product
        product.stock -= 1;

        // Increment total products sold
        totalProductsSold += 1;
    }

    // Function to close the store (only owner can close the store)
    function closeStore() public onlyOwner {
        storeStatus = StoreStatus.Closed;
    }

    // Function to reopen the store (only owner can reopen the store)
    function openStore() public onlyOwner {
        storeStatus = StoreStatus.Open;
    }

    // Function to get the number of products in the store
    function getProductCount() public view returns (uint) {
        return products.length;
    }

    // Function to get the details of a specific product
    function getProduct(uint _productIndex) public view returns (string memory name, uint price, uint stock) {
        require(_productIndex < products.length, "Product does not exist");
        Product storage product = products[_productIndex];
        return (product.name, product.price, product.stock);
    }
}
