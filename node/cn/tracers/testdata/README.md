# Making testdata

1. Create a transaction on chain. You can use any network - local private net, testnet, mainnet. Remember its txhash.
2. Run the following command to generate a JSON testdata. Necessary context data will be compiled in the JSON.

    ```sh
    export RPC=http://localhost:8551
    ./makeTest.sh fastCallTracer 0x6bcce4a683a1e81168e7ab05c3b4fa7d17a1cb97a70ef5c666a14e4603615b0c call_tracer/my_test.json
    ./makeTest.sh prestateTracer 0x6bcce4a683a1e81168e7ab05c3b4fa7d17a1cb97a70ef5c666a14e4603615b0c prestate_tracer/my_test.json
    ```

3. Inspect the JSON file.
4. Edit the JSON file.
    - Delete the last line saying `undefined`.
    - Edit the `"_comment"` field in the JSON as needed.

