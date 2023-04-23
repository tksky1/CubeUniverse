import {
    Title,
    Box,
    Text,
    Grid,
    useMantineTheme,
    Group,
    ThemeIcon,
    LoadingOverlay
} from "@mantine/core"
import { IconExclamationCircle } from "@tabler/icons-react";
import dynamic from "next/dynamic";
import MyCard from "../MyCard";

const ResponsivePie = dynamic(() => import('../re-exports').then(module => module.ResponsivePie), {
    loading: () => <LoadingOverlay visible={true} overlayBlur={2} />,
    ssr: false
})

export default function Capacity() {
    let theme = useMantineTheme();
    return (
        <Box p={"md"}>
            <Group
                align={"center"}
                spacing={"xs"}
                mb={"xs"}
            >
                <Title order={6} color={theme.colors.gray[6]}>
                    Capacity
                </Title>
                <ThemeIcon size={"xs"}>
                    <IconExclamationCircle />
                </ThemeIcon>
            </Group>
            <Grid>
                <Grid.Col span={3}>
                    <MyCard title="Raw Capacity">
                        <Box sx={theme => ({
                            width: "100%",
                            height: "300px"
                        })}>
                            <MyPie></MyPie>
                        </Box>
                    </MyCard>
                </Grid.Col>
                <Grid.Col span={3}>
                    <MyCard title="Hosts">
                        <Text align="center" weight={"bold"}>
                            3 total
                        </Text>
                    </MyCard>
                </Grid.Col>
                <Grid.Col span={3}>
                    <MyCard title="Monitors">
                        <Text weight={"bold"}>
                            3 (quorum 0, 1, 2)
                        </Text>
                    </MyCard>
                </Grid.Col>
                <Grid.Col span={3}>
                    <MyCard title="OSDs">
                        <>
                            <Text weight={"bold"}>
                                3 total
                            </Text>
                            <Text weight={"bold"}>
                                3 up, 3 in
                            </Text>
                        </>
                    </MyCard>
                </Grid.Col>
                <Grid.Col span={3}>
                    <MyCard title="Managers">
                        <>
                            <Text weight={"bold"}>
                                1 active
                            </Text>
                            <Text weight={"bold"}>
                                1 standby
                            </Text>
                        </>
                    </MyCard>
                </Grid.Col>
                <Grid.Col span={3}>
                    <MyCard title="Object Gateways">
                        <Text weight={"bold"}>
                            0 total
                        </Text>
                    </MyCard>
                </Grid.Col>
                <Grid.Col span={3}>
                    <MyCard title="Metadata Servers">
                        <Text weight={"bold"}>
                            no filesystems
                        </Text>
                    </MyCard>
                </Grid.Col>
                <Grid.Col span={3}>
                    <MyCard title="iSCSC Gateways">
                        <>
                            <Text weight={"bold"}>
                                0 total
                            </Text>
                            <Text weight={"bold"}>
                                0 up, 0 dowm
                            </Text>
                        </>
                    </MyCard>
                </Grid.Col>
            </Grid>
        </Box>
    )
}

let data = [
    {
        "id": "c",
        "label": "c",
        "value": 427,
        "color": "hsl(358, 70%, 50%)"
    },
    {
        "id": "ruby",
        "label": "ruby",
        "value": 92,
        "color": "hsl(146, 70%, 50%)"
    }
]

function MyPie() {
    return (
        <ResponsivePie
            data={data}
            margin={{ bottom: 80, top: 40, left: 40, right: 40 }}
            innerRadius={0.7}
            padAngle={0.7}
            cornerRadius={3}
            activeOuterRadiusOffset={8}
            borderWidth={2}
            borderColor={{
                from: 'color',
                modifiers: [
                    [
                        'darker',
                        0.2
                    ]
                ]
            }}
            arcLinkLabelsSkipAngle={10}
            arcLinkLabelsTextColor="#333333"
            arcLinkLabelsThickness={2}
            arcLinkLabelsColor={{ from: 'color' }}
            arcLabelsSkipAngle={10}
            arcLabelsTextColor={{
                from: 'color',
                modifiers: [
                    [
                        'darker',
                        2
                    ]
                ]
            }}
            defs={[
                {
                    id: 'dots',
                    type: 'patternDots',
                    background: 'inherit',
                    color: 'rgba(255, 255, 255, 0.3)',
                    size: 4,
                    padding: 1,
                    stagger: true
                },
                {
                    id: 'lines',
                    type: 'patternLines',
                    background: 'inherit',
                    color: 'rgba(255, 255, 255, 0.3)',
                    rotation: -45,
                    lineWidth: 6,
                    spacing: 10
                }
            ]}
            fill={[
                {
                    match: {
                        id: 'ruby'
                    },
                    id: 'dots'
                },
                {
                    match: {
                        id: 'c'
                    },
                    id: 'dots'
                },
            ]}
            legends={[
                {
                    anchor: 'bottom',
                    direction: 'row',
                    justify: false,
                    translateX: 0,
                    translateY: 56,
                    itemsSpacing: 0,
                    itemWidth: 100,
                    itemHeight: 18,
                    itemTextColor: '#999',
                    itemDirection: 'left-to-right',
                    itemOpacity: 1,
                    symbolSize: 18,
                    symbolShape: 'circle',
                    effects: [
                        {
                            on: 'hover',
                            style: {
                                itemTextColor: '#000'
                            }
                        }
                    ]
                }
            ]}
        />
    )
}
