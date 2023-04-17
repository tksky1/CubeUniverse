import { commonLayout } from "@/utils/commonLayout";
import { useContext } from "react"
import { DataContext } from "@/components/DataContext";
import { Box } from "@mantine/core";
import { MantineReactTable } from 'mantine-react-table';
import type { MRT_ColumnDef } from 'mantine-react-table';

export default function CephPool() {
    return (
        <Box p={"xl"}>
            <MyTable />
        </Box>
    )
}

interface CephPool {
    Name: string,
    Replica: number,
    PG: number,
    CreateTime: string
}

let columns: MRT_ColumnDef<CephPool>[] = [
    {
        header: "存储池",
        accessorKey: "Name"
    },
    {
        header: "副本个数",
        accessorKey: "Replica"
    },
    {
        header: "对象组",
        accessorKey: "PG"
    },
    {
        header: "创建时间",
        accessorKey: "CreateTime"
    }
]

function MyTable() {
    let data = useContext(DataContext);
    let pools: CephPool[] = data && (data as any).CephPools ? (data as any).CephPools : [];
    return (
        <MantineReactTable
            columns={columns}
            data={pools ? pools : []}
            enableRowSelection
            enableColumnOrdering
            enableGlobalFilter
            enableBottomToolbar
            enableClickToCopy
            mantineTableProps={{
                sx: theme => ({
                    backgroundColor: "#f7f8fb"
                })
            }}
            mantineTableBodyProps={{
                sx: theme => ({
                    backgroundColor: "#f7f8fb"
                })
            }}
        />
    )
}

CephPool.getLayout = commonLayout;