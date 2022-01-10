# Write your own docker

This is a project to write a docker by yourself.

> PLEASE RUN PROJECT WITH LINUX (like ubuntu).
>
> Or you may failed 😥

## Step 1：Initial container namespace

![carbon](https://typora-photos.oss-cn-shenzhen.aliyuncs.com/carbon.png)

## Step 2：Initial container resource limit

This step we initial resource limit through Cgroup. Please be carefully while testing cpu limit, you may got different answers with me -- the cpu percent is 100%, this is because your computer is multi-cores. If you want to try, buy an one-core cloudy service machine is ok. 🤓

![carbon-2](https://typora-photos.oss-cn-shenzhen.aliyuncs.com/carbon-2.png)

## Step 3：Use busybox container

This step we use a small image - busybox, and through pivot_root help our docker to run this container. Trust me, it will amazing you! 🤩

![carbon-3](https://typora-photos.oss-cn-shenzhen.aliyuncs.com/carbon-3.png)

## Little tips

If you code with Goland but run project in virtual machine, you can exec this command.

scp -rp GolandProjects/writeDocker ubuntu@192.168.64.13:/home/ubuntu

![carbon1](https://typora-photos.oss-cn-shenzhen.aliyuncs.com/carbon1.png)