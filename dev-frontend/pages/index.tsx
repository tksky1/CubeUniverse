import { Badge, Box, Grid, LoadingOverlay, Paper, Title, Stack, Center, Text, useMantineTheme, Loader } from "@mantine/core";
import { linearGradientDef } from '@nivo/core'
import dynamic from "next/dynamic";
import { useState, useEffect, useRef } from "react"
import { commonLayout } from "@/utils/commonLayout";
import { useContext } from "react"
import MyCard from "@/components/MyCard";
import useWebSocket from "react-use-websocket";
import { DataContext } from "@/components/DataContext";

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

export default function Home() {
  let data = useContext(DataContext);
  return (
    <Box p={"lg"}>
      <Grid grow gutter={5} gutterXs="md" gutterMd="xl" gutterXl={50} align={"stretch"}>
        <Grid.Col
          span={4}
          sx={theme => ({
            height: "200px"
          })}
        >
          <MyOprLine />
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
            height: "200px"
          })}
        >
          <MyByteLine />
        </Grid.Col>
        <Grid.Col span={4}>
          <ObjectPie />
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

function MyOprLine() {
  let dataLength = 20;
  let { lastMessage } = useWebSocket("ws://control-backend.cubeuniverse.svc.cluster.local:30401/api/storage/pvcws");
  let [data, setData] = useState<{ id: string, data: { x: number, y: number }[] }[]>([
    {
      id: "读操作数",
      data: []
    },
    {
      id: "写操作数",
      data: []
    }
  ]);
  useEffect(() => {
    if (lastMessage) {
      let msgData = JSON.parse(lastMessage.data);
      if (msgData.CephPerformance) {
        console.log(msgData.CephPerformance.ReadOperationsPerSec, msgData.CephPerformance.WriteOperationPerSec)
        setData(lastData => {
          let read = lastData[0].data;
          let write = lastData[0].data;
          if (read.length === 0) {
            read = [{
              x: 1,
              y: msgData.CephPerformance.ReadOperationsPerSec
            }];
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
          }
          return [
            {
              id: "读操作数",
              data: read
            },
            {
              id: "写操作数",
              data: write
            }
          ]
        })
      }
    }
  }, [lastMessage]);
  return (
    <Paper
      shadow={"md"}
      sx={theme => ({
        backgroundColor: "#f7f8fb",
        height: "100%",
        width: "100%",
      })}
    >
      <ResponsiveLine
        margin={{ left: 60, bottom: 30, right: 120, top: 30 }}
        data={data}
        curve="natural"
        colors={{ scheme: 'accent' }}
        lineWidth={1}
        enableArea
        areaOpacity={0.5}
        axisBottom={null}
        pointSize={5}
        defs={[
          linearGradientDef('gradientA', [
            { offset: 50, color: 'inherit' },
            { offset: 100, color: 'inherit', opacity: 0 },
          ]),
        ]}
        fill={[{ match: '*', id: 'gradientA' }]}
        legends={[
          {
            anchor: 'bottom-right',
            direction: 'column',
            justify: false,
            translateX: 100,
            translateY: 0,
            itemsSpacing: 0,
            itemDirection: 'left-to-right',
            itemWidth: 80,
            itemHeight: 20,
            itemOpacity: 0.75,
            symbolSize: 12,
            symbolShape: 'circle',
            symbolBorderColor: 'rgba(0, 0, 0, .5)',
            effects: [
              {
                on: 'hover',
                style: {
                  itemBackground: 'rgba(0, 0, 0, .03)',
                  itemOpacity: 1
                }
              }
            ]
          }
        ]}
      />
    </Paper>
  )
}

function MyByteLine() {
  let dataLength = 20;
  let { lastMessage } = useWebSocket("ws://control-backend.cubeuniverse.svc.cluster.local:30401/api/storage/pvcws");
  let [data, setData] = useState<{ id: string, data: { x: number, y: number }[] }[]>([
    {
      id: "读比特数",
      data: []
    },
    {
      id: "写比特数",
      data: []
    }
  ]);
  useEffect(() => {
    if (lastMessage) {
      let msgData = JSON.parse(lastMessage.data);
      if (msgData.CephPerformance) {
        console.log(msgData.CephPerformance.ReadOperationsPerSec, msgData.CephPerformance.WriteOperationPerSec)
        setData(lastData => {
          let read = lastData[0].data;
          let write = lastData[0].data;
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
          return [
            {
              id: "读比特数",
              data: read
            },
            {
              id: "写比特数",
              data: write
            }
          ]
        })
      }
    }
  }, [lastMessage]);
  return (
    <Paper
      shadow={"md"}
      sx={theme => ({
        backgroundColor: "#f7f8fb",
        height: "100%",
        width: "100%",
      })}
    >
      <ResponsiveLine
        margin={{ left: 60, bottom: 30, right: 120, top: 30 }}
        data={data}
        curve="natural"
        colors={{ scheme: 'set3' }}
        lineWidth={1}
        enableArea
        areaOpacity={0.5}
        axisBottom={null}
        pointSize={5}
        defs={[
          linearGradientDef('gradientA', [
            { offset: 50, color: 'inherit' },
            { offset: 100, color: 'inherit', opacity: 0 },
          ]),
        ]}
        fill={[{ match: '*', id: 'gradientA' }]}
        legends={[
          {
            anchor: 'bottom-right',
            direction: 'column',
            justify: false,
            translateX: 100,
            translateY: 0,
            itemsSpacing: 0,
            itemDirection: 'left-to-right',
            itemWidth: 80,
            itemHeight: 20,
            itemOpacity: 0.75,
            symbolSize: 12,
            symbolShape: 'circle',
            symbolBorderColor: 'rgba(0, 0, 0, .5)',
            effects: [
              {
                on: 'hover',
                style: {
                  itemBackground: 'rgba(0, 0, 0, .03)',
                  itemOpacity: 1
                }
              }
            ]
          }
        ]}
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

function ObjectPie() {
  let data = useContext(DataContext);
  let theme = useMantineTheme();
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
              id: "丢失",
              value: (data as any).CephPerformance.ObjectNotFoundNum
            },
            {
              id: "未归置",
              value: (data as any).CephPerformance.ObjectMisplacedNum
            },
            {
              id: "降级",
              value: (data as any).CephPerformance.ObjectDegradedNum
            },
            {
              id: "正常",
              value: (data as any).CephPerformance.ObjectReplicatedNum
                - (data as any).CephPerformance.ObjectNotFoundNum
                - (data as any).CephPerformance.ObjectMisplacedNum
                - (data as any).CephPerformance.ObjectDegradedNum
            }
          ]
        }
        margin={{ top: 0, right: 160, bottom: 20, left: 60 }}
        innerRadius={0.6}
        padAngle={3}
        cornerRadius={3}
        borderWidth={1}
        arcLinkLabelsStraightLength={0}
        activeOuterRadiusOffset={4}
        arcLinkLabelsSkipAngle={5}
        arcLabelsSkipAngle={5}
        tooltip={({ datum }) => {
          return (
            <Paper
              p={"md"}
              shadow={"xl"}
              bg={theme.colors.dark[4]}
            >
              <Center>
                <Text color={datum.color}>
                  {datum.id}: {datum.value}
                </Text>
              </Center>
            </Paper>
          )
        }}
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

Home.getLayout = commonLayout;