import { Badge, Box, Grid, Paper, Title, Stack, Center, Loader } from "@mantine/core";
import { commonLayout } from "@/utils/commonLayout";
import { useContext } from "react"
import MyCard from "@/components/MyCard";
import { DataContext } from "@/components/DataContext";
import EChartsReact from "echarts-for-react";
import { graphic } from "echarts";
import { bytesRData, bytesWData, oprRData, oprWData } from "@/components/DataProvider";
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

export default function Home() {
    let data = useContext(DataContext);
    data && console.log((data as any).CephPerformance);
    console.log(bytesRData.value, bytesWData.value, oprRData.value, ...oprWData.value);
    return (
        <Box p={"lg"}>
            <Grid grow gutter={5} gutterXs="md" gutterMd="xl" gutterXl={50} align={"stretch"}>
                <Grid.Col
                    span={6}
                    sx={theme => ({
                        height: "400px"
                    })}
                >
                    <NewOprLine />
                </Grid.Col>
                <Grid.Col
                    span={6}
                    sx={theme => ({
                        height: "400px"
                    })}
                >
                    <NewBytePie />
                </Grid.Col>
                <Grid.Col span={6}
                    sx={theme => ({
                        height: "400px"
                    })}>
                    <NewByteLine />
                    {/* <NewNewOprLine /> */}
                </Grid.Col>
                <Grid.Col span={6}>
                    <NewObjectPie />
                </Grid.Col>
                <Grid.Col
                    span={4}
                >
                    <Paper
                        shadow={"xl"}
                        sx={theme => ({
                            height: "100%"
                        })}
                    >
                        <MyCard title="状态总览">
                            <Badge
                                color={!(data && (data as any).CephPerformance)
                                    ? "blue"
                                    : (data as any).CephPerformance.HealthStatus === "HEALTH_WARN"
                                        ? "yellow"
                                        : (data as any).CephPerformance.HealthStatus === "HEALTH_OK"
                                            ? "green"
                                            : (data as any).CephPerformance.HealthStatus === "HEALTH_ERROR"
                                                ? "red"
                                                : "blue"}
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
                >
                    <Paper
                        shadow={"xl"}
                        sx={theme => ({
                            height: "100%"
                        })}
                    >
                        <MyCard title="健康状况详情">
                            <Stack spacing={"xs"}>
                                <Title order={3} align="center">
                                    {data && (data as any).CephPerformance
                                        ? (data as any).CephPerformance.HealthStatusDetailed.map((t: any) => <>{t}</>)
                                        : "暂无数据"}
                                </Title>
                            </Stack>
                        </MyCard>
                    </Paper>
                </Grid.Col>
                <Grid.Col
                    span={4}
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
                    span={4}
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
                    span={4}
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
                    span={4}
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
                    span={4}
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
                    span={4}
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
                    span={4}
                >
                    <Paper
                        shadow={"xl"}
                        sx={theme => ({
                            height: "100%"
                        })}
                    >
                        <MyCard title="就绪OSD数量">
                            <Stack spacing={"xs"}>
                                <Title order={3}>
                                    {data && (data as any).CephPerformance
                                        ? `${(data as any).CephPerformance.OSDNotReadyNum}`
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
                style={{ height: '100%', width: '100%' }}
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

const data = [
    {
        name: "Page A",
        uv: 4000,
        pv: 2400,
        amt: 2400
    },
    {
        name: "Page B",
        uv: 3000,
        pv: 1398,
        amt: 2210
    },
    {
        name: "Page C",
        uv: 2000,
        pv: 9800,
        amt: 2290
    },
    {
        name: "Page D",
        uv: 2780,
        pv: 3908,
        amt: 2000
    },
    {
        name: "Page E",
        uv: 1890,
        pv: 4800,
        amt: 2181
    },
    {
        name: "Page F",
        uv: 2390,
        pv: 3800,
        amt: 2500
    },
    {
        name: "Page G",
        uv: 3490,
        pv: 4300,
        amt: 2100
    }
];

function NewNewOprLine() {
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
            <ResponsiveContainer>
                <AreaChart
                    data={oprRData.value.map((x, index) => {
                        let w = oprWData.value;
                        return {
                            x: x.x,
                            write: w[index].y,
                            read: x.y
                        }
                    })}
                    margin={{
                        top: 10,
                        right: 30,
                        left: 0,
                        bottom: 0
                    }}
                >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="x" />
                    <YAxis />
                    <Tooltip />
                    <Area
                        type="monotone"
                        dataKey="read"
                        stackId="1"
                        stroke="#8884d8"
                        fill="#8884d8"
                    />
                    <Area
                        type="monotone"
                        dataKey="write"
                        stackId="1"
                        stroke="#82ca9d"
                        fill="#82ca9d"
                    />
                </AreaChart>
            </ResponsiveContainer>
        </Paper>
    )
}


function NewByteLine() {
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
                            areaStyle: {
                                color: new graphic.LinearGradient(0, 0, 0, 1, [
                                    {
                                        offset: 0,
                                        color: 'rgb(255, 158, 68)'
                                    },
                                    {
                                        offset: 1,
                                        color: 'rgb(255, 70, 131)'
                                    }
                                ])
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
                                color: 'rgb(106, 90, 205)'
                            },
                            areaStyle: {
                                color: new graphic.LinearGradient(0, 0, 0, 1, [
                                    {
                                        offset: 0,
                                        color: 'rgb(132, 112, 255)'
                                    },
                                    {
                                        offset: 1,
                                        color: 'rgb(106, 90, 205)'
                                    }
                                ])
                            },
                            data: bytesWData.value.map(x => x.y)
                        }
                    ]
                }}
            />
        </Paper>
    )
}

function NewBytePie() {
    let data = useContext(DataContext);
    return data && (data as any).CephPerformance ? (
        <Paper
            shadow={"md"}
            sx={theme => ({
                backgroundColor: "#f7f8fb",
                height: "100%",
                width: "100%"
            })}
        >
            <EChartsReact
                style={{ height: '100%', width: '100%' }}
                option={{
                    tooltip: {
                        formatter: '{a} <br/>{b} : {c}%'
                    },
                    series: [
                        {
                            name: 'Pressure',
                            type: 'gauge',
                            max: Math.trunc(((data as any).CephPerformance.TotalBytes / 1024 / 1024 / 1024) * 10) / 10,
                            // max: 4,
                            progress: {
                                show: true
                            },
                            detail: {
                                valueAnimation: true,
                                formatter: '{value} GB'
                            },
                            data: [
                                {
                                    value: Math.trunc((((data as any).CephPerformance.TotalBytes - (data as any).CephPerformance.TotalUsedBytes) / 1024 / 1024 / 1024) * 10) / 10,
                                    name: '可用容量'
                                }
                            ]
                        }
                    ]
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