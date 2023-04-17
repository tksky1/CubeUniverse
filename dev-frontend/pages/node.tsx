import { commonLayout } from "@/utils/commonLayout";
import { useContext } from "react"
import { DataContext } from "@/components/DataContext";
import dynamic from "next/dynamic";
import { Box, Badge, Space, LoadingOverlay } from "@mantine/core";
import { MantineReactTable } from 'mantine-react-table';
import type { MRT_ColumnDef } from 'mantine-react-table';

const ResponsiveLine = dynamic(
    () => import('../components/re-exports').then(module => module.ResponsiveLine),
    {
        loading: () => <LoadingOverlay visible={true} overlayBlur={2} />,
        ssr: false
    }
)

export default function Node() {
    return (
        <Box p={"xl"}>
            <HostTable />
            <Space h="md" />
            <MonitorTable />
        </Box>
    )
}

interface Host {
    Hostname: string,
    Services: string[]
}

let hostColumns: MRT_ColumnDef<Host>[] = [
    {
        header: "节点名称",
        accessorKey: "Hostname"
    },
    {
        header: "服务列表",
        accessorKey: "Services",
        Cell: ({ cell }) => {
            return (
                <Box>
                    {cell
                        .getValue<string[]>()
                        .map((s, index) => <Badge key={`service-${index}`} mx={"xs"} variant={"outline"}>
                            {s}
                        </Badge>)}
                </Box>
            )
        }
    }
]

function HostTable() {
    let data = useContext(DataContext);
    let hosts: Host[] = data && (data as any).CephHosts ? (data as any).CephHosts : [];
    return (
        <MantineReactTable
            columns={hostColumns}
            data={hosts ? hosts : []}
            enableRowSelection
            enableColumnOrdering
            enableGlobalFilter
            enableBottomToolbar
        />
    )
}

interface CephMonitor {
    Name: string,
    Rank: number,
    Address: string,
    OpenSessions: number[]
}

let monitorColumns: MRT_ColumnDef<CephMonitor>[] = [
    {
        header: "监控节点",
        accessorKey: "Name"
    },
    {
        header: "权重",
        accessorKey: "Rank"
    },
    {
        header: "地址",
        accessorKey: "Address"
    },
    {
        header: "活跃的监控事务数",
        accessorKey: "OpenSessions",
        Cell: ({ cell }) => {
            let data = cell.getValue<number[]>().slice(0, 10);
            for (let i = data.length; i < 10; i++) {
                data[i] = 0;
            }
            let min = data.slice(1).reduce((prev, current) => current > prev ? prev : current, data[0]);
            let max = data.slice(1).reduce((prev, current) => current > prev ? current : prev, data[0]);
            let tickValues = Array(max - min + 1);
            for (let i = 0; i < tickValues.length; i++) {
                tickValues[i] = min + i;
            }
            return (
                <Box
                    w={300}
                    h={40 * (max - min + 1)}
                >
                    <ResponsiveLine
                        data={[
                            {
                                id: 1,
                                data: data.map((val, index) => ({ x: index, y: val }))
                            }
                        ]}
                        margin={{ left: 25, bottom: 30, right: 30, top: 30 }}
                        enableGridX={false}
                        enableGridY={false}
                        axisBottom={null}
                        yScale={{
                            type: "linear",
                            min: "auto",
                            max: "auto",
                            stacked: true
                        }}
                        axisLeft={{
                            tickValues
                        }}
                        yFormat=" >-.0d"
                    />
                </Box>
            )
        }
    }
]

function MonitorTable() {
    let data = useContext(DataContext);
    let monitors: CephMonitor[] = data && (data as any).inQuorumMonitors ? (data as any).inQuorumMonitors : [];
    return (
        <MantineReactTable
            columns={monitorColumns}
            data={monitors ? monitors : []}
            enableRowSelection
            enableColumnOrdering
            enableGlobalFilter
            enableBottomToolbar
        />
    )
}


Node.getLayout = commonLayout;