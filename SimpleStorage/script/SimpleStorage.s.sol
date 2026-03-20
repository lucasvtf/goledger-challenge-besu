// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.27;

import {Script, console} from "forge-std/Script.sol";
import {SimpleStorage} from "../src/SimpleStorage.sol";

contract SimpleStorageScript is Script {
    function run() public returns (SimpleStorage) {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");

        vm.startBroadcast(deployerPrivateKey);

        SimpleStorage simpleStorage = new SimpleStorage();

        vm.stopBroadcast();

        console.log("SimpleStorage deployed at:", address(simpleStorage));

        return simpleStorage;
    }
}
