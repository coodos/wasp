// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmlib from "wasmlib"
import * as sc from "./index";

export class DonateCall {
    func: wasmlib.ScFunc = new wasmlib.ScFunc(sc.HScName, sc.HFuncDonate);
    params: sc.MutableDonateParams = new sc.MutableDonateParams();
}

export class DonateContext {
    params: sc.ImmutableDonateParams = new sc.ImmutableDonateParams();
    state: sc.MutableDonateWithFeedbackState = new sc.MutableDonateWithFeedbackState();
}

export class WithdrawCall {
    func: wasmlib.ScFunc = new wasmlib.ScFunc(sc.HScName, sc.HFuncWithdraw);
    params: sc.MutableWithdrawParams = new sc.MutableWithdrawParams();
}

export class WithdrawContext {
    params: sc.ImmutableWithdrawParams = new sc.ImmutableWithdrawParams();
    state: sc.MutableDonateWithFeedbackState = new sc.MutableDonateWithFeedbackState();
}

export class DonationCall {
    func: wasmlib.ScView = new wasmlib.ScView(sc.HScName, sc.HViewDonation);
    params: sc.MutableDonationParams = new sc.MutableDonationParams();
    results: sc.ImmutableDonationResults = new sc.ImmutableDonationResults();
}

export class DonationContext {
    params: sc.ImmutableDonationParams = new sc.ImmutableDonationParams();
    results: sc.MutableDonationResults = new sc.MutableDonationResults();
    state: sc.ImmutableDonateWithFeedbackState = new sc.ImmutableDonateWithFeedbackState();
}

export class DonationInfoCall {
    func: wasmlib.ScView = new wasmlib.ScView(sc.HScName, sc.HViewDonationInfo);
    results: sc.ImmutableDonationInfoResults = new sc.ImmutableDonationInfoResults();
}

export class DonationInfoContext {
    results: sc.MutableDonationInfoResults = new sc.MutableDonationInfoResults();
    state: sc.ImmutableDonateWithFeedbackState = new sc.ImmutableDonateWithFeedbackState();
}

export class ScFuncs {

    static donate(ctx: wasmlib.ScFuncCallContext): DonateCall {
        let f = new DonateCall();
        f.func.setPtrs(f.params, null);
        return f;
    }

    static withdraw(ctx: wasmlib.ScFuncCallContext): WithdrawCall {
        let f = new WithdrawCall();
        f.func.setPtrs(f.params, null);
        return f;
    }

    static donation(ctx: wasmlib.ScViewCallContext): DonationCall {
        let f = new DonationCall();
        f.func.setPtrs(f.params, f.results);
        return f;
    }

    static donationInfo(ctx: wasmlib.ScViewCallContext): DonationInfoCall {
        let f = new DonationInfoCall();
        f.func.setPtrs(null, f.results);
        return f;
    }
}
