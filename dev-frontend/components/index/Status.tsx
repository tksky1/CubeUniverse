import { 
    Title, 
    Box, 
    Text, 
    Grid, 
    useMantineTheme, 
    Group, 
    ThemeIcon,
    Badge,
} from "@mantine/core"
import { IconExclamationCircle, IconAlertTriangle } from "@tabler/icons-react";
import MyCard from "../MyCard";

export default function Status() {
    let theme = useMantineTheme();
    return (
        <Box p={"md"}>
            <Group
                align={"center"}
                spacing={"xs"}
                mb={"xs"}
            >
                <Title order={6} color={theme.colors.gray[6]}>
                    Status
                </Title>
                <ThemeIcon size={"xs"}>
                    <IconExclamationCircle/>
                </ThemeIcon>
            </Group>
            <Grid grow>
                <Grid.Col span={3}>
                    <MyCard title="Cluster Status">
                        <Badge 
                            variant={"outline"} 
                            radius="xs" 
                            size={"xl"}
                            color="yellow"
                        >
                            <Group spacing={1} noWrap>
                                <Text>
                                    HEALTH_WARM
                                </Text>
                                <ThemeIcon 
                                    color="yellow"
                                    variant={"outline"}
                                    sx={theme => ({
                                        borderStyle: "none"
                                    })}
                                >
                                    <IconAlertTriangle/>
                                </ThemeIcon>
                            </Group>
                        </Badge>
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