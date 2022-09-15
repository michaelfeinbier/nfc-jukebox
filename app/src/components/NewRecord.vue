<script setup>
import { reactive, ref } from 'vue';

import { StreamBarcodeReader } from "vue-barcode-reader";
const barcode = ref(null)
//const state = reactive()

function onDecode(detectedBarcode) {
    barcode.value = detectedBarcode
    console.log(detectedBarcode)
}

</script>
<template>
    <div class="form border p-4">
        <h4>Scan record</h4>
        <div class="mb-3">
            <StreamBarcodeReader @decode="onDecode"></StreamBarcodeReader>
        </div>
        <div class="mb-3">
            <label for="" class="form-label">Detected Barcode:</label>
            <input type="text" ref="barcode" class="form-control" />
        </div>
    </div>
</template>
