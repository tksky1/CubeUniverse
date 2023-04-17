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
    Group
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
                        <BlockPVCTable pvcs={pvcs} afterSubmit={updatePvc} />
                    </Box>
                )
                : (
                    <Center
                        h={"100%"}
                    >
                        <Stack>
                            <Title order={3}>
                                Not Open
                            </Title>
                            <Button onClick={() => {
                                openStorage("block")
                                    .then(r => r.json())
                                    .then(r => r.code === 200
                                        ? Promise.resolve(r) : Promise.reject(r))
                                    .then(r => noti("Success", "Success", "Open Storage"))
                                    .catch(r => noti("Error", "Error", r.msg))
                                    .finally(() => updateStorageOpen());
                            }}>
                                Open Me
                            </Button>
                        </Stack>
                    </Center>
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
                            label="Name"
                            {...form.getInputProps("name")} />
                        <TextInput
                            withAsterisk
                            label="NameSpace"
                            {...form.getInputProps("namespace")} />
                        <NumberInput
                            min={1}
                            label="Volume"
                            {...form.getInputProps("volume")} />
                        <Switch
                            label="Autoscale"
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
                        label="Volume"
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