// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.27;

import {Test} from "forge-std/Test.sol";
import {SimpleStorage} from "../src/SimpleStorage.sol";

contract SimpleStorageTest is Test {
    SimpleStorage public simpleStorage;

    function setUp() public {
        simpleStorage = new SimpleStorage();
    }

    function test_InitialValueIsZero() public view {
        assertEq(simpleStorage.get(), 0);
    }

    function test_Set() public {
        simpleStorage.set(42);
        assertEq(simpleStorage.get(), 42);
    }

    function test_SetOverwritesPreviousValue() public {
        simpleStorage.set(100);
        simpleStorage.set(999);
        assertEq(simpleStorage.get(), 999);
    }

    function testFuzz_Set(uint256 x) public {
        simpleStorage.set(x);
        assertEq(simpleStorage.get(), x);
    }
}
