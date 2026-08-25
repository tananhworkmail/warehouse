<template>
    <tr>
        <!-- 主標題 -->
        <th
            v-for="PrimaryHeader in props.headers"
            :key="PrimaryHeader.title"
            :colspan="
                PrimaryHeader.children ? PrimaryHeader.children.length : 1
            "
            :rowspan="PrimaryHeader.children ? 1 : 2"
            style="
                white-space: nowrap;
                text-align: center;
                background-color: #f5f5f5;
                font-size: 14px;
                border: 1px solid lightgray;
            "
        >
            {{ PrimaryHeader.title }}
        </th>
    </tr>
    <tr>
        <!-- 子標題 -->
        <th
            v-for="(ChildrenHeader, i1) in getChildrenHeader(props.headers)"
            :key="i1"
            style="
                white-space: nowrap;
                text-align: center;
                background-color: #f5f5f5;
                font-size: 14px;
                padding: 0;
                border: 1px solid lightgray;
            "
        >           
            {{ ChildrenHeader.title }}
        </th>
    </tr>
    <tr>
        <!-- 次標題 -->
        <th
            v-for="(SubHeader, i3) in getEngSubHeader(props.headers)"
            :key="i3"
            style="
                white-space: nowrap;
                text-align: center;
                background-color: #f5f5f5;
                color: black;
                font-weight: 600;
                font-family: 'Times New Roman', Times, serif;
                font-size: 14px;
                border: 1px solid lightgray;
            "
        >
            {{ SubHeader.title }}
        </th>
    </tr>
    <tr>
        <!-- 次標題 -->
        <th
            v-for="(SubHeader, i2) in getVietSubHeader(props.headers)"
            :key="i2"
            style="
                white-space: nowrap;
                text-align: center;
                background-color: black;
                color: cyan;
                font-weight: 600;
                font-family: 'Times New Roman', Times, serif;
                font-size: 14px;
                border: 1px solid lightgray;
            "
        >
            {{ SubHeader.title }}
        </th>
    </tr>
</template>

<script setup>
const props = defineProps(['headers'])
const getChildrenHeader = (headers) => {
    let result = []
    headers
        .filter((i) => i.children)
        .forEach((v) => {
            result = result.concat(v.children)
        })
    return result
}
const getVietSubHeader = (headers) => {
    let result = []
    headers
        .filter((i) => i.subvietheader)
        .forEach((v) => {
            result = result.concat(v.subvietheader)
        })
    return result
}
const getEngSubHeader = (headers) => {
    let result = []
    headers
        .filter((i) => i.subengheader)
        .forEach((v) => {
            result = result.concat(v.subengheader)
        })
    return result
}
</script>
<style></style>
