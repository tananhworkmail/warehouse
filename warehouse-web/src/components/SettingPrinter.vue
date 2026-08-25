<template>
    <div class="centered-container">
        <v-dialog v-model="dialog" max-width="600">
            <template v-slot:activator="{ props: activatorProps }">
                <v-btn v-bind="activatorProps" class="printer-ip-button">
                    <v-icon left icon="$ip" class="printer-ip-icon"></v-icon>
                    PRINTER IP
                </v-btn>
            </template>

            <v-card prepend-icon="$ip" title="SETTING PRINTER IP">
                <v-card-text>
                    <v-row dense>
                        <v-col cols="3">
                            <v-text-field
                                density="compact"
                                v-model="ipPart1"
                                readonly
                                variant="outlined"
                            ></v-text-field>
                        </v-col>
                        <v-col cols="3">
                            <v-text-field
                                density="compact"
                                v-model="ipPart2"
                                variant="outlined"
                                readonly
                            ></v-text-field>
                        </v-col>
                        <v-col cols="3">
                            <v-text-field
                                density="compact"
                                v-model="ipPart3"
                                variant="outlined"
                                required
                                :rules="[validateIPPart]"
                                maxlength="3"
                                clearable
                                placeholder="XXX"
                            ></v-text-field>
                        </v-col>
                        <v-col cols="3">
                            <v-text-field
                                density="compact"
                                v-model="ipPart4"
                                variant="outlined"
                                required
                                :rules="[validateIPPart]"
                                maxlength="3"
                                clearable
                                placeholder="XXX"
                            ></v-text-field>
                        </v-col>
                    </v-row>

                    <small
                        v-if="settingMessage"
                        class="text-caption"
                        :style="{
                            color: settingMessage.includes('Success')
                                ? 'green'
                                : 'red',
                        }"
                    >
                        * {{ settingMessage }}
                    </small>

                    <!-- Loading Animation -->
                    <div v-if="isLoading" class="loading-spinner">
                        <v-progress-linear
                            color="cyan"
                            indeterminate
                        ></v-progress-linear>
                    </div>
                </v-card-text>

                <v-divider></v-divider>

                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn
                        text="Clear IP Setting"
                        variant="outlined"
                        @click="clearIPSetting"
                    ></v-btn>
                    <v-btn
                        color="primary"
                        text="Save"
                        variant="outlined"
                        @click="savePrinterIP"
                    ></v-btn>

                    <v-btn
                        text="Close"
                        variant="plain"
                        @click="dialog = false"
                    ></v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { usePrinter } from '@/hooks/usePrinter'

const { setPrinterIP, clearPrinterIP } = usePrinter()
const dialog = ref(false)
const ipPart1 = ref('192')
const ipPart2 = ref('168')
const ipPart3 = ref(
    localStorage.getItem('printerIP')
        ? localStorage.getItem('printerIP').split('.')[2]
        : '',
)
const ipPart4 = ref(
    localStorage.getItem('printerIP')
        ? localStorage.getItem('printerIP').split('.')[3]
        : '',
)

const settingMessage = ref('')
const isLoading = ref(false) // Loading state

// 驗證 IP 部分是否有效（0-255）
const validateIPPart = (value) => {
    const partPattern = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
    return partPattern.test(value) || 'Invalid IP segment (0-255)'
}

const PrintSuccessMSG = async (ip) => {
    const zplCommand = `^XA^FO50,50^A0N,50,50^FDsetting successful to IP ${ip}^FS^XZ`

    try {
        const response = await fetch(`http://${ip}/pstprnt`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded', // ZPL 的 Content-Type
            },
            body: zplCommand,
        })

        if (response.ok) {
            console.log('印刷命令已成功發送到印表機:', ip)
            return true
        } else {
            console.error('印表機連線失敗:', response.status)
            return 'Printer connection failed: ' + response.status
        }
    } catch (error) {
        console.error('印表機連線失敗:', error)
        return 'Printer connection failed: ' + error.message
    }
}

const clearIPSetting = () => {
    ipPart3.value = ''
    ipPart4.value = ''
    clearPrinterIP()
}
const savePrinterIP = async () => {
    settingMessage.value = '' // Clear any previous error messages
    isLoading.value = true // Start loading animation

    // Validate the IP parts
    const isValidPart3 = validateIPPart(ipPart3.value)
    const isValidPart4 = validateIPPart(ipPart4.value)

    if (isValidPart3 === true && isValidPart4 === true) {
        const fullIP = `${ipPart1.value}.${ipPart2.value}.${ipPart3.value}.${ipPart4.value}`

        // Attempt to set the printer IP and wait for the result
        const success = await setPrinterIP(fullIP) // Add await here

        if (success) {
            localStorage.setItem('printerIP', fullIP) // Save the full IP to localStorage
            setTimeout(() => {
                isLoading.value = false // Stop loading animation
                settingMessage.value = 'Success to set Printer IP.' // Handle failure to set printer IP
                PrintSuccessMSG(fullIP)
            }, 2000)
        } else {
            setTimeout(() => {
                isLoading.value = false // Stop loading animation
                settingMessage.value = 'Failed to set Printer IP.' // Handle failure to set printer IP
            }, 2000)
        }
    } else {
        // Handle invalid IP parts
        settingMessage.value =
            isValidPart3 !== true ? isValidPart3 : isValidPart4
        console.error('Invalid IP Address:', ipPart3.value, ipPart4.value)
    }

    // Simulate loading for 5 seconds
}
</script>

<style scoped>
.centered-container {
    display: flex;
    align-items: center;
    justify-content: center;
}

.printer-ip-button {
    background-image: linear-gradient(to left, #0ba360 0%, #3cba92 100%);
    color: white;
    font-weight: 500 !important;
    display: flex;
}

.printer-ip-icon {
    color: white;
    font-size: 20px;
    margin-right: 8px;
}

.printer-ip-title {
    display: flex;
    align-items: center;
    color: green;
}

.v-btn {
    font-weight: 500 !important;
}

.loading-spinner {
    display: flex;
    justify-content: center;
    margin-top: 10px; /* Optional: Adjust margin as needed */
}
</style>
