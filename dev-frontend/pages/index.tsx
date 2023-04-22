import { Badge, Box, Grid, LoadingOverlay, Paper, Title, Stack, Center, Text, useMantineTheme, Loader } from "@mantine/core";
import { linearGradientDef } from '@nivo/core'
import dynamic from "next/dynamic";
import { useState, useEffect, useRef } from "react"
import { commonLayout } from "@/utils/commonLayout";
import { useContext } from "react"
import MyCard from "@/components/MyCard";
import { DataContext } from "@/components/DataContext";
import { atom } from "signia";
import EChartsReact from "echarts-for-react";
import echarts from "echarts";

const ResponsiveLine = dynamic(
    () => import('../components/re-exports').then(module => module.ResponsiveLine),
    {
        loading: () => <LoadingOverlay visible={true} overlayBlur={2} />,
        ssr: false
    }
)

const ResponsivePie = dynamic(
    () => import('../components/re-exports').then(module => module.ResponsivePie),
    {
        loading: () => <LoadingOverlay visible={true} overlayBlur={2} />,
        ssr: false
    }
)

let bytesWData = atom<{ x: number, y: number }[]>("bytesData", []);
let bytesRData = atom<{ x: number, y: number }[]>("bytesData", []);
let oprWData = atom<{ x: number, y: number }[]>("oprData", []);
let oprRData = atom<{ x: number, y: number }[]>("oprData", []);

export default function Home() {
    let data = useContext(DataContext);
    return (
        <Box p={"lg"}>
            <Grid grow gutter={5} gutterXs="md" gutterMd="xl" gutterXl={50} align={"stretch"}>
                <Grid.Col
                    span={4}
                    sx={theme => ({
                        height: "350px"
                    })}
                >
                    <NewOprLine />
                </Grid.Col>
                <Grid.Col span={4}>
                    <ByteUsagePie />
                </Grid.Col>
                <Grid.Col
                    span={3}
                >
                    <Paper
                        shadow={"xl"}
                        sx={theme => ({
                            height: "100%"
                        })}
                    >
                        <MyCard title="状态总览">
                            <Badge
                                color={"yellow"}
                                size="xl"
                                radius="xs"
                            >
                                {(data && (data as any).CephPerformance) ? (data as any).CephPerformance.HealthStatus : "暂无数据"}
                            </Badge>
                        </MyCard>
                    </Paper>
                </Grid.Col>
                <Grid.Col
                    span={4}
                    sx={theme => ({
                        height: "350px"
                    })}
                >
                    <NewByteLine />
                </Grid.Col>
                <Grid.Col span={4}>
                    <NewObjectPie />
                </Grid.Col>
                <Grid.Col
                    span={3}
                >
                    <Paper
                        shadow={"xl"}
                        sx={theme => ({
                            height: "100%"
                        })}
                    >
                        <MyCard title="存储池总数">
                            <Stack spacing={"xs"}>
                                <Title order={2}>
                                    {data && (data as any).CephPerformance
                                        ? `${(data as any).CephPerformance.PoolNum}`
                                        : "暂无数据"}
                                </Title>
                            </Stack>
                        </MyCard>
                    </Paper>
                </Grid.Col>
                <Grid.Col
                    span={2}
                >
                    <Paper
                        shadow={"xl"}
                        sx={theme => ({
                            height: "100%"
                        })}
                    >
                        <MyCard title="存储节点总数">
                            <Title order={3}>
                                {data && (data as any).CephPerformance ? (data as any).CephPerformance.HostNum : "暂无数据"}
                            </Title>
                        </MyCard>
                    </Paper>
                </Grid.Col>
                <Grid.Col
                    span={2}
                >
                    <Paper
                        shadow={"xl"}
                        sx={theme => ({
                            height: "100%"
                        })}
                    >
                        <MyCard title="监控节点总数">
                            <Title order={3}>
                                {data && (data as any).CephPerformance
                                    ? `${(data as any).CephPerformance.MonitorNum}`
                                    : "暂无数据"}
                            </Title>
                        </MyCard>
                    </Paper>
                </Grid.Col>
                <Grid.Col
                    span={2}
                >
                    <Paper
                        shadow={"xl"}
                        sx={theme => ({
                            height: "100%"
                        })}
                    >
                        <MyCard title="每秒恢复数据流量">
                            <Stack spacing={"xs"}>
                                <Title order={3}>
                                    {data && (data as any).CephPerformance
                                        ? `${(data as any).CephPerformance.RecoveringBytesPerSec} bytes/s`
                                        : "暂无数据"}
                                </Title>
                            </Stack>
                        </MyCard>
                    </Paper>
                </Grid.Col>

                <Grid.Col
                    span={2}
                >
                    <Paper
                        shadow={"xl"}
                        sx={theme => ({
                            height: "100%"
                        })}
                    >
                        <MyCard title="独立对象总数">
                            <Stack spacing={"xs"}>
                                <Title order={3}>
                                    {data && (data as any).CephPerformance
                                        ? `${(data as any).CephPerformance.ObjectNum}`
                                        : "暂无数据"}
                                </Title>
                            </Stack>
                        </MyCard>
                    </Paper>
                </Grid.Col>
                <Grid.Col
                    span={2}
                >
                    <Paper
                        shadow={"xl"}
                        sx={theme => ({
                            height: "100%"
                        })}
                    >
                        <MyCard title="每秒恢复数据流量">
                            <Stack spacing={"xs"}>
                                <Title order={3}>
                                    {data && (data as any).CephPerformance
                                        ? `${(data as any).CephPerformance.RecoveringBytesPerSec} bytes/s`
                                        : "暂无数据"}
                                </Title>
                            </Stack>
                        </MyCard>
                    </Paper>
                </Grid.Col>
            </Grid>
        </Box>
    )
}

function NewOprLine() {
    let dataLength = 20;
    let fatherData = useContext(DataContext);
    useEffect(() => {
        if (fatherData) {
            let msgData = fatherData as any;
            if (msgData.CephPerformance) {
                console.log(msgData.CephPerformance.ReadOperationsPerSec, msgData.CephPerformance.WriteOperationPerSec)
                let read = oprRData.value;
                let write = oprWData.value;
                if (read.length === 0) {
                    read = [{
                        x: 1,
                        y: msgData.CephPerformance.ReadOperationsPerSec
                    }]
                } else {
                    let shouldShift = read.length >= dataLength;
                    read.push({
                        x: read[read.length - 1].x + 1,
                        y: msgData.CephPerformance.ReadOperationsPerSec
                    });
                    shouldShift && read.shift();
                }
                if (write.length === 0) {
                    write = [{
                        x: 1,
                        y: msgData.CephPerformance.WriteOperationPerSec
                    }];
                } else {
                    let shouldShift = write.length >= dataLength;
                    write.push({
                        x: write[write.length - 1].x + 1,
                        y: msgData.CephPerformance.WriteOperationPerSec
                    });
                    shouldShift && write.shift();
                };
                oprRData.set(read);
                oprWData.set(write);
            }
        }
    }, [fatherData]);
    return (
        <Paper
            shadow={"md"}
            sx={theme => ({
                backgroundColor: "#f7f8fb",
                height: "100%",
                width: "100%",
            })}
            p={"md"}
        >
            <EChartsReact
                option={{
                    title: {
                        text: 'CubeUniverse I/O 操作数'
                    },
                    tooltip: {},
                    xAxis: {
                        data: oprRData.value.map(x => x.x),
                        splitLine: {
                            show: false
                        }
                    },
                    yAxis: {},
                    series: [
                        {
                            name: '读操作数',
                            type: 'bar',
                            data: oprRData.value.map(x => x.y),
                            emphasis: {
                                focus: 'series'
                            },
                            animationDelay: function (idx: number) {
                                return idx * 10;
                            }
                        },
                        {
                            name: '写操作数',
                            type: 'bar',
                            data: oprWData.value.map(x => x.y),
                            emphasis: {
                                focus: 'series'
                            },
                            animationDelay: function (idx: number) {
                                return idx * 10 + 100;
                            }
                        }
                    ],
                    animationEasing: 'elasticOut',
                    animationDelayUpdate: function (idx: number) {
                        return idx * 5;
                    }
                }}
            />
        </Paper>
    )
}

function NewByteLine() {
    let dataLength = 20;
    let fatherData = useContext(DataContext);
    useEffect(() => {
        if (fatherData) {
            let msgData = fatherData as any;
            if (msgData.CephPerformance) {
                console.log(msgData.CephPerformance.ReadBytesPerSec, msgData.CephPerformance.ReadBytesPerSec)
                let read = bytesRData.value;
                let write = bytesWData.value;
                if (read.length === 0) {
                    read = [{
                        x: 1,
                        y: msgData.CephPerformance.ReadBytesPerSec
                    }];
                } else {
                    let shouldShift = read.length >= dataLength;
                    read.push({
                        x: read[read.length - 1].x + 1,
                        y: msgData.CephPerformance.ReadBytesPerSec
                    });
                    shouldShift && read.shift();
                }
                if (write.length === 0) {
                    write = [{
                        x: 1,
                        y: msgData.CephPerformance.WriteBytesPerSec
                    }];
                } else {
                    let shouldShift = write.length >= dataLength;
                    write.push({
                        x: write[write.length - 1].x + 1,
                        y: msgData.CephPerformance.WriteBytesPerSec
                    });
                    shouldShift && write.shift();
                }
                bytesRData.set(read);
                bytesWData.set(write);
            }
        }
    }, [fatherData]);
    return (
        <Paper
            shadow={"md"}
            sx={theme => ({
                backgroundColor: "#f7f8fb",
                height: "100%",
                width: "100%",
            })}
            p={"md"}
        >
            <EChartsReact
                option={{
                    tooltip: {
                        trigger: 'axis',
                        position: function (pt: number[]) {
                            return [pt[0], '10%'];
                        },
                        axisPointer: {
                            type: 'cross'
                        }
                    },
                    title: {
                        text: 'CubeUniverse I/O 流量'
                    },
                    xAxis: {
                        type: 'category',
                        boundaryGap: false,
                        data: bytesWData.value.map(x => x.x)
                    },
                    yAxis: {
                        type: 'value',
                        boundaryGap: [0, '30%']
                    },
                    dataZoom: [
                        {
                            type: 'inside',
                            start: 0,
                            end: 100
                        },
                        {
                            start: 0,
                            end: 100
                        }
                    ],
                    series: [
                        {
                            name: '读流量(字节)',
                            type: 'line',
                            showSymbol: false,
                            smooth: true,
                            sampling: 'lttb',
                            itemStyle: {
                                color: 'rgb(255, 70, 131)'
                            },
                            data: bytesRData.value.map(x => x.y)
                        },
                        {
                            name: '写流量(字节)',
                            type: 'line',
                            showSymbol: false,
                            smooth: true,
                            sampling: 'lttb',
                            itemStyle: {
                                color: 'rgb(106 90 205)'
                            },
                            data: bytesWData.value.map(x => x.y)
                        }
                    ]
                }}
            />
        </Paper>
    )
}

function ByteUsagePie() {
    let data = useContext(DataContext);
    return data && (data as any).CephPerformance ? (
        <Paper
            shadow={"md"}
            sx={theme => ({
                backgroundColor: "#f7f8fb",
                height: "100%",
                width: "100%",
            })}
        >
            <ResponsivePie
                data={
                    [
                        {
                            id: "已用容量",
                            value: (data as any).CephPerformance.TotalUsedBytes
                        },
                        {
                            id: "可用容量",
                            value: (data as any).CephPerformance.TotalBytes - (data as any).CephPerformance.TotalUsedBytes
                        }
                    ]
                }
                margin={{ top: 0, right: 160, bottom: 20, left: 60 }}
                innerRadius={0.6}
                padAngle={3}
                cornerRadius={3}
                borderWidth={1}
                arcLinkLabelsStraightLength={0}
                activeOuterRadiusOffset={8}
                arcLinkLabelsSkipAngle={5}
                arcLabelsSkipAngle={5}
                legends={[
                    {
                        anchor: 'right',
                        direction: 'column',
                        justify: false,
                        translateX: 160,
                        translateY: 0,
                        itemWidth: 100,
                        itemHeight: 20,
                        itemsSpacing: 10,
                        symbolSize: 20,
                        itemDirection: 'left-to-right'
                    }
                ]}
            />
        </Paper>
    ) : (
        <Paper
            shadow={"md"}
            sx={theme => ({
                backgroundColor: "#f7f8fb",
                height: "100%",
                width: "100%",
            })}
        >
            <Center
                h={"100%"}
            >
                <Loader />
            </Center>
        </Paper>
    )
}

function NewObjectPie() {
    let data = useContext(DataContext);
    return data && (data as any).CephPerformance ? (
        <Paper
            shadow={"md"}
            p={"xs"}
            sx={theme => ({
                backgroundColor: "#f7f8fb",
                height: "100%",
                width: "100%",
            })}
        >
            <EChartsReact
                option={{
                    title: {
                        text: "对象组状态",
                        subtext: "对象存储服务内对象",
                    },
                    tooltip: {
                        show: true,
                        trigger: "item",
                        formatter: "{b}: {c} ({d}%)",
                    },
                    legend: {
                        orient: "horizontal",
                        bottom: "0%",
                        data: ["<10w", "10w-50w", "50w-100w", "100w-500w", ">500w"],
                    },
                    series: [
                        {
                            type: "pie",
                            selectedMode: "single",
                            radius: ["25%", "58%"],
                            color: ["#86D560", "#AF89D6", "#59ADF3", "#FF999A", "#FFCC67"],

                            label: {
                                normal: {
                                    position: "inner",
                                    formatter: "{d}%",

                                    textStyle: {
                                        color: "#fff",
                                        fontWeight: "bold",
                                        fontSize: 14,
                                    },
                                },
                            },
                            labelLine: {
                                normal: {
                                    show: false,
                                },
                            },
                            data: [
                                {
                                    name: "丢失",
                                    value: (data as any).CephPerformance.ObjectNotFoundNum
                                },
                                {
                                    name: "未归置",
                                    value: (data as any).CephPerformance.ObjectMisplacedNum
                                },
                                {
                                    name: "降级",
                                    value: (data as any).CephPerformance.ObjectDegradedNum
                                },
                                {
                                    name: "正常",
                                    value: (data as any).CephPerformance.ObjectReplicatedNum
                                        - (data as any).CephPerformance.ObjectNotFoundNum
                                        - (data as any).CephPerformance.ObjectMisplacedNum
                                        - (data as any).CephPerformance.ObjectDegradedNum
                                }
                            ],
                        },
                        {
                            type: "pie",
                            radius: ["58%", "83%"],
                            itemStyle: {
                                normal: {
                                    color: "#F2F2F2",
                                },
                                emphasis: {
                                    color: "#ADADAD",
                                },
                            },
                            label: {
                                normal: {
                                    position: "inner",
                                    formatter: "{c}",
                                    textStyle: {
                                        color: "#777777",
                                        fontWeight: "bold",
                                        fontSize: 14,
                                    },
                                },
                            },
                            data: [
                                {
                                    name: "丢失",
                                    value: (data as any).CephPerformance.ObjectNotFoundNum
                                },
                                {
                                    name: "未归置",
                                    value: (data as any).CephPerformance.ObjectMisplacedNum
                                },
                                {
                                    name: "降级",
                                    value: (data as any).CephPerformance.ObjectDegradedNum
                                },
                                {
                                    name: "正常",
                                    value: (data as any).CephPerformance.ObjectReplicatedNum
                                        - (data as any).CephPerformance.ObjectNotFoundNum
                                        - (data as any).CephPerformance.ObjectMisplacedNum
                                        - (data as any).CephPerformance.ObjectDegradedNum
                                }
                            ],
                        },
                    ],
                }}
            />
        </Paper>
    ) : (
        <Paper
            shadow={"md"}
            sx={theme => ({
                backgroundColor: "#f7f8fb",
                height: "100%",
                width: "100%",
            })}
        >
            <Center
                h={"100%"}
            >
                <Loader />
            </Center>
        </Paper>
    )
}

Home.getLayout = commonLayout;