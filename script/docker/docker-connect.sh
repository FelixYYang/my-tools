#!/bin/bash

# docker 容器快速连接脚本

# 设置颜色变量
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BLUE='\033[0;34m'

docker_connect() {
    # 检查 Docker 是否在运行
    if ! docker info >/dev/null 2>&1; then
        echo -e "${RED}Docker 服务未运行！${NC}"
        return 1
    fi

    # 如果没有提供参数，显示容器列表供选择
    if [ $# -eq 0 ]; then
        echo -e "${BLUE}运行中的容器列表：${NC}"
        containers=($(docker ps --format "{{.Names}}"))
        
        if [ ${#containers[@]} -eq 0 ]; then
            echo -e "${RED}没有正在运行的容器！${NC}"
            return 1
        fi

        # 显示容器列表
        for i in "${!containers[@]}"; do
            container_id=$(docker ps --filter "name=${containers[$i]}" --format "{{.ID}}")
            container_image=$(docker ps --filter "name=${containers[$i]}" --format "{{.Image}}")
            echo -e "$((i+1)): ${GREEN}${containers[$i]}${NC} (ID: ${container_id:0:12}, Image: $container_image)"
        done

        # 获取用户选择
        echo -e "\n${BLUE}请选择容器编号 [1-${#containers[@]}]:${NC} "
        read -r choice

        # 验证输入
        if ! [[ "$choice" =~ ^[0-9]+$ ]] || [ "$choice" -lt 1 ] || [ "$choice" -gt "${#containers[@]}" ]; then
            echo -e "${RED}无效的选择！${NC}"
            return 1
        fi

        selected_container="${containers[$((choice-1))]}"
    else
        selected_container="$1"
    fi

    # 验证容器是否存在且运行
    if ! docker ps --format "{{.Names}}" | grep -q "^${selected_container}$"; then
        echo -e "${RED}容器 '$selected_container' 未运行或不存在！${NC}"
        return 1
    fi

    # 检测容器中可用的shell
    if docker exec "$selected_container" which bash >/dev/null 2>&1; then
        shell="bash"
    else
        shell="sh"
    fi

    echo -e "${GREEN}正在连接到容器 '$selected_container' 使用 $shell...${NC}"
    docker exec -it "$selected_container" "$shell"
}

# 创建命令别名
alias dc='docker_connect'

# 添加命令补全功能
_docker_connect_completion() {
    local curr_word="${COMP_WORDS[COMP_CWORD]}"
    local containers=$(docker ps --format "{{.Names}}")
    COMPREPLY=($(compgen -W "$containers" -- "$curr_word"))
}

complete -F _docker_connect_completion docker_connect
complete -F _docker_connect_completion dc
