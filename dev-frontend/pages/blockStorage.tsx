import { commonLayout } from "@/utils/commonLayout"
import {
    Button,
    Center,
    Modal,
    NumberInput,
    Stack,
    Switch,
    TextInput,
    Title,
    Tooltip,
    ActionIcon,
    Loader,
    Group,
    Paper,
    Text,
    Highlight,
    useMantineTheme
} from "@mantine/core";
import { useEffect, useState, useCallback } from "react";
import { checkStorageOpen, createPvc, deletePvc, getPvcInfo, openStorage, updatePvcVolume } from "@/apis";
import { notifications } from "@mantine/notifications"
import { Box } from "@mantine/core";
import { MantineReactTable } from 'mantine-react-table';
import type { MRT_ColumnDef } from 'mantine-react-table';
import { IconX, IconCheck, IconEdit, IconTrash } from "@tabler/icons-react";
import { useDisclosure } from "@mantine/hooks";
import { useForm } from "@mantine/form";
import { v4 } from "uuid";
import { noti } from "@/utils/noti";
import blockStorage from "../public/block.png"
import Image from "next/image";

interface BlockPVC {
    name: string,
    namespace: string,
    time: string,
    volume: string
}

export default function BlockStorage() {
    let [storageOpen, setStorageOpen] = useState<"loading" | boolean>("loading");
    let [pvcs, setPvcs] = useState<BlockPVC[]>([]);
    let updateStorageOpen = useCallback(() => {
        checkStorageOpen("block")
            .then(e => e.json())
            .then(e => {
                if (e.msg && e.msg[0] === "p") {
                    setStorageOpen(true);
                } else {
                    setStorageOpen(false);
                }
            });
    }, []);
    let updatePvc = useCallback(() => (
        getPvcInfo("block")
            .then(r => r.json())
            .then(r => setPvcs(Object.values(r.data)))
            .catch(r => {
                notifications.show({
                    id: v4(),
                    withCloseButton: true,
                    autoClose: 5000,
                    title: "Error",
                    message: "Something went wrong",
                    color: 'red',
                    icon: <IconX />,
                    loading: false,
                    variant: "outline"
                });
            })), []);
    useEffect(() => {
        updateStorageOpen();
        updatePvc();
    }, []);
    return (
        storageOpen === "loading"
            ? (
                <Center h="100%">
                    <Loader size={"lg"} />
                </Center>
            )
            : storageOpen
                ? (
                    <Box
                        p={"md"}
                        sx={theme => ({
                        })}>
                        <Stack>
                            <BlockStorageInfo />
                            <BlockPVCTable pvcs={pvcs} afterSubmit={updatePvc} />
                        </Stack>
                    </Box>
                )
                : (
                    <BlockNotOpen openIt={() => {
                        openStorage("block")
                            .then(r => r.json())
                            .then(r => r.code === 200
                                ? Promise.resolve(r) : Promise.reject(r))
                            .then(r => noti("Success", "Success", "Open Storage"))
                            .catch(r => noti("Error", "Error", r.msg))
                            .finally(() => updateStorageOpen());
                    }} />
                )
    )
}

BlockStorage.getLayout = commonLayout;

let hostColumns: MRT_ColumnDef<BlockPVC>[] = [
    {
        header: "存储卷声明",
        accessorKey: "name"
    },
    {
        header: " 命名空间",
        accessorKey: "namespace",
    },
    {
        header: "存储卷容量",
        accessorKey: "volume"
    },
    {
        header: "创建时间",
        accessorKey: "time"
    }
]

interface BlockPVCTableInterface {
    pvcs: BlockPVC[],
    afterSubmit: () => Promise<void>
}

function BlockPVCTable({ pvcs, afterSubmit }: BlockPVCTableInterface) {
    let form = useForm({
        initialValues: {
            name: "",
            namespace: "",
            volume: 1,
            autoscale: false
        },
        validate: {
            name: val => val.length === 0 ? "Empty Name" : null,
            namespace: val => val.length === 0 ? "Empty Namespace" : null,
            volume: val => (typeof val === "number" && val >= 1) ? null : "Invalid Value"
        }
    });
    let [opened, { open, close }] = useDisclosure(false);
    let [editOpened, { close: editClose, open: editOpen }] = useDisclosure(false);
    let [deleteOpened, { close: deleteClose, open: deleteOpen }] = useDisclosure(false);
    let [currentRow, setCurrentRow] = useState<null | { name: string, namespace: string }>(null);
    useEffect(() => {
        form.reset();
    }, [opened]);
    return (
        <>
            <MantineReactTable
                columns={hostColumns}
                data={pvcs}
                enableColumnOrdering
                enableGlobalFilter
                enableBottomToolbar
                renderTopToolbarCustomActions={() => {
                    return (
                        <Button
                            variant="gradient"
                            gradient={{ from: 'teal', to: 'lime', deg: 105 }}
                            onClick={open}>
                            创建PVC
                        </Button>)
                }}
                onEditingRowSave={() => { }}
                enableEditing
                renderRowActions={({ row, table }) => (
                    <Box sx={{ display: 'flex', gap: '16px' }}>
                        <Tooltip withArrow position="left" label="编辑">
                            <ActionIcon
                                onClick={() => {
                                    setCurrentRow({
                                        name: row.original.name,
                                        namespace: row.original.namespace
                                    });
                                    editOpen();
                                }}
                            >
                                <IconEdit />
                            </ActionIcon>
                        </Tooltip>
                        <Tooltip withArrow position="right" label="删除">
                            <ActionIcon
                                color="red"
                                onClick={() => {
                                    setCurrentRow({
                                        name: row.original.name,
                                        namespace: row.original.namespace
                                    });
                                    deleteOpen();
                                }}>
                                <IconTrash />
                            </ActionIcon>
                        </Tooltip>
                    </Box>
                )}
            />
            <Modal
                opened={opened}
                onClose={close}
                centered
                title={<Title order={4}>新建PVC</Title>}
                radius={"lg"}
                withCloseButton={false}>
                <form onSubmit={form.onSubmit(values => {
                    createPvc({
                        name: values.name,
                        "X-type": "block",
                        autoscale: values.autoscale ? "true" : "false",
                        namespace: values.namespace,
                        volume: `${values.volume}`
                    }).then(e => e.json())
                        .then((e) => {
                            return e.code === 200
                                ? Promise.resolve(e)
                                : Promise.reject(e);
                        })
                        .then(() => afterSubmit().then(() => {
                            close();
                            notifications.show({
                                id: v4(),
                                withCloseButton: true,
                                autoClose: 5000,
                                title: "Success",
                                message: "Create new PVC",
                                color: "green",
                                icon: <IconCheck />,
                                loading: false,
                            });
                        })).catch(e => {
                            close();
                            notifications.show({
                                id: v4(),
                                withCloseButton: true,
                                autoClose: 5000,
                                title: "Error",
                                message: e.msg,
                                color: "red",
                                icon: <IconX />,
                                loading: false,
                            });
                        });
                })}>
                    <Stack>
                        <TextInput
                            withAsterisk
                            label="存储卷声明"
                            {...form.getInputProps("name")} />
                        <TextInput
                            withAsterisk
                            label="命名空间"
                            {...form.getInputProps("namespace")} />
                        <NumberInput
                            min={1}
                            label="存储卷容量"
                            {...form.getInputProps("volume")} />
                        <Switch
                            label="自动扩容"
                            {...form.getInputProps("autoscale")} />
                        <Button
                            variant="gradient"
                            gradient={{ from: 'teal', to: 'blue', deg: 105 }}
                            type="submit">
                            确定
                        </Button>
                    </Stack>
                </form>
            </Modal>
            <EditModal
                opened={editOpened}
                close={editClose}
                name={currentRow ? currentRow.name : ""}
                namespace={currentRow ? currentRow.namespace : ""}
                afterSubmit={afterSubmit}
            />
            <DeleteModal
                afterSubmit={afterSubmit}
                opened={deleteOpened}
                close={deleteClose}
                name={currentRow ? currentRow.name : ""}
                namespace={currentRow ? currentRow.namespace : ""}
            />
        </>
    )
}

interface EditModalInterface {
    name: string,
    namespace: string,
    afterSubmit: () => Promise<void>,
    opened: boolean,
    close: VoidFunction
}

function EditModal({ name, namespace, afterSubmit, opened, close }: EditModalInterface) {
    let form = useForm({
        initialValues: {
            volume: 1,
        },
        validate: {
            volume: val => (typeof val === "number" && val >= 1) ? null : "Invalid Value"
        }
    });
    useEffect(() => {
        form.reset();
    }, [opened]);
    return (
        <Modal
            opened={opened}
            onClose={close}
            centered
            title={<Title order={4}>编辑</Title>}
            size={"xs"}
            radius={"lg"}
            withCloseButton={false}>
            <form onSubmit={form.onSubmit(values => {
                updatePvcVolume({
                    name: name,
                    "X-type": "block",
                    namespace: namespace,
                    volume: `${values.volume}`
                }).then(e => e.json())
                    .then((e) => {
                        return e.code === 200
                            ? Promise.resolve(e)
                            : Promise.reject(e);
                    })
                    .then(() => afterSubmit().then(() => {
                        close();
                        notifications.show({
                            id: v4(),
                            withCloseButton: true,
                            autoClose: 5000,
                            title: "Success",
                            message: "Update PVC colume",
                            color: "green",
                            icon: <IconCheck />,
                            loading: false,
                        });
                    })).catch(e => {
                        close();
                        notifications.show({
                            id: v4(),
                            withCloseButton: true,
                            autoClose: 5000,
                            title: "Error",
                            message: e.msg,
                            color: "red",
                            icon: <IconX />,
                            loading: false,
                        });
                    });
            })}>
                <Stack>
                    <NumberInput
                        size="md"
                        min={1}
                        label="存储卷容量"
                        {...form.getInputProps("volume")} />
                    <Button
                        variant="gradient"
                        gradient={{ from: 'teal', to: 'blue', deg: 105 }}
                        type="submit">
                        确定
                    </Button>
                </Stack>
            </form>
        </Modal>
    )
}

interface DeleteModalInterface {
    name: string,
    namespace: string,
    afterSubmit: () => Promise<void>,
    opened: boolean,
    close: VoidFunction
}

function DeleteModal({ name, namespace, afterSubmit, opened, close }: DeleteModalInterface) {
    return (
        <Modal
            opened={opened}
            onClose={close}
            centered
            size={"auto"}
            radius={"md"}
            withCloseButton={false}>
            <Stack>
                <Group position="apart" w="100%">
                    <Title order={5}>确定要删除此项PVC吗？</Title>
                    <Group spacing={"xs"}>
                        <Button size="xs" color="gray" onClick={close}>取消</Button>
                        <Button size="xs" color="red" onClick={() => {
                            deletePvc({ name, namespace, "X-type": "block" })
                                .then(e => e.json())
                                .then(e => e.code === 200 ? Promise.resolve(e) : Promise.reject(e))
                                .then(e => {
                                    notifications.show({
                                        id: v4(),
                                        withCloseButton: true,
                                        autoClose: 5000,
                                        title: "Success",
                                        message: "PVC deleted",
                                        color: "green",
                                        icon: <IconCheck />,
                                        loading: false,
                                    });
                                })
                                .catch(e => {
                                    notifications.show({
                                        id: v4(),
                                        withCloseButton: true,
                                        autoClose: 5000,
                                        title: "Error",
                                        message: e.msg,
                                        color: 'red',
                                        icon: <IconX />,
                                        loading: false,
                                        variant: "outline"
                                    });
                                })
                                .finally(() => {
                                    close();
                                    afterSubmit();
                                })
                        }}>删除</Button>
                    </Group>
                </Group>
            </Stack>
        </Modal>
    )
}
function BlockNotOpen({ openIt }: { openIt: any }) {
    let theme = useMantineTheme();
    return (
        <Box h={"100%"}>
            <Center h={"100%"}>
                <Paper
                    radius={"xl"}
                    shadow="xl"
                    w={"60%"}
                    p={"xl"}>
                    <Stack>
                        <Box>
                            <Image
                                src={blockStorage}
                                height={100}
                                style={{
                                    float: "left",
                                    margin: 10
                                }}
                                alt="BlockStorage" />
                            <Highlight
                                size={"md"}
                                highlight={['块存储', '数据分块', "一对一"]}
                                highlightStyles={(theme) => ({
                                    backgroundImage: theme.fn.linearGradient(45, theme.colors.grape[5], theme.colors.orange[5]),
                                    fontWeight: 700,
                                    WebkitBackgroundClip: 'text',
                                    WebkitTextFillColor: 'transparent',
                                })}
                            >
                                块存储是一种直接基于底层物理设备直接提供存储的技术。它将数据分块并建立快速访问索引，避免大量的元数据传输从而提高存储性能并降低延迟。因此，基于块存储的存储卷能提供比基于文件存储的存储卷更快的访问速度，但不支持多个应用同时访问。在 Kubernetes中以持久化卷的形式提供，常用于高存储性能要求的应用，仅支持一对一绑定。在控制面板创建PVC后，为集群应用pod绑定该PVC即可挂载。
                            </Highlight>
                        </Box>
                        <Center>
                            <Button
                                onClick={() => openIt()}
                                variant="gradient"
                                gradient={{ from: theme.colors.grape[5], to: theme.colors.orange[5], deg: 45 }}>开启块存储</Button>
                        </Center>
                    </Stack>
                </Paper>
            </Center>
        </Box>
    )
}

function BlockStorageInfo() {
    let theme = useMantineTheme();
    return (
        <Center>
            <Paper
                p={"xl"}
                sx={theme => ({
                    borderRadius: 1,
                    boxShadow: "0px 2px 5px 1px rgba(196,74,24,0.75)",
                    border: "1px solid #e67c2c"
                })}>
                <Stack>
                    <Box mr={100} ml={50} mt={20}>
                        <Image
                            height={100}
                            src={blockStorage}
                            style={{
                                float: "left",
                                marginRight: 30,
                            }}
                            alt="FileStorage" />
                        <Title
                            align="left"
                            order={3}
                            mb={"md"}
                        >块存储</Title>
                        <Highlight
                            size={"md"}
                            highlight={['块存储', '数据分块', "一对一"]}
                            highlightStyles={(theme) => ({
                                backgroundImage: theme.fn.linearGradient(45, theme.colors.grape[5], theme.colors.orange[5]),
                                fontWeight: 700,
                                WebkitBackgroundClip: 'text',
                                WebkitTextFillColor: 'transparent',
                            })}
                        >
                            块存储是一种直接基于底层物理设备直接提供存储的技术。它将数据分块并建立快速访问索引，避免大量的元数据传输从而提高存储性能并降低延迟。因此，基于块存储的存储卷能提供比基于文件存储的存储卷更快的访问速度，但不支持多个应用同时访问。在 Kubernetes中以持久化卷的形式提供，常用于高存储性能要求的应用，仅支持一对一绑定。在控制面板创建PVC后，为集群应用pod绑定该PVC即可挂载。
                        </Highlight>
                    </Box>
                </Stack>
            </Paper>
        </Center>

    )
}