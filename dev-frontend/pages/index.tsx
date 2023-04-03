import { Box, Grid, LoadingOverlay, Paper } from "@mantine/core";
import { linearGradientDef } from '@nivo/core'
import dynamic from "next/dynamic";
import { useState, useEffect, useRef } from "react"
import { commonLayout } from "@/utils/commonLayout";
import Status from "@/components/index/Status";
import Capacity from "@/components/index/Capacity";

const ResponsiveLine = dynamic(
  () => import('../components/re-exports').then(module => module.ResponsiveLine),
  {
    loading: () => <LoadingOverlay visible={true} overlayBlur={2} />,
    ssr: false
  })

export default function Home() {
  return (
    <Box p={"lg"}>
      <Grid gutter={5} gutterXs="md" gutterMd="xl" gutterXl={50}>
        <Grid.Col
          span={8}
          sx={theme => ({
            height: "300px"
          })}
        >
          <MyLine />
        </Grid.Col>
        <Grid.Col
          span={8}
          sx={theme => ({
            height: "300px"
          })}
        >
          <MyLine />
        </Grid.Col>
      </Grid>
      {/* <Status/>
      <Capacity/> */}
    </Box>
  )
}

function MyLine() {
  let dataLength = 20;
  let [data, setData] = useState<{ id: string, data: { x: number, y: number }[] }[]>([
    {
      id: "joe",
      data: []
    },
    {
      id: "tim",
      data: []
    }
  ]);
  let id = useRef<ReturnType<typeof setInterval> | null>(null);
  useEffect(() => {
    id.current = setInterval(() => {
      setData(origin => {
        origin.forEach(x => updateData(x.data));
        return origin.slice();
      })
    }, 1000);
    return () => {
      id.current && clearInterval(id.current);
    }
  }, []);
  function updateData(data: { x: number, y: number }[]) {
    let dataItemLength = data.length;
    if (dataItemLength !== dataLength) {
      data.push({
        x: dataItemLength,
        y: Math.random() * dataLength
      });
      return;
    }
    let shifted = data.shift();
    data.push({ x: shifted!.x + dataItemLength, y: Math.random() * dataLength });
  }
  return (
    <Paper
      shadow={"xl"}
      sx={theme => ({
        height: "100%",
        width: "100%",
      })}
    >
      <ResponsiveLine
        margin={{ left: 60, bottom: 30, right: 60, top: 30 }}
        data={data}
        curve="natural"
        colors={{ scheme: 'dark2' }}
        lineWidth={1}
        enableArea
        areaOpacity={0.5}
        // enablePoints={false}
        pointSize={5}
        defs={[
          linearGradientDef('gradientA', [
            { offset: 50, color: 'inherit' },
            { offset: 100, color: 'inherit', opacity: 0 },
          ]),
        ]}
        fill={[{ match: '*', id: 'gradientA' }]}
      />
    </Paper>
  )
}

Home.getLayout = commonLayout;