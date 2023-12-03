//
// Created by Guilherme on 01/12/2023.
//
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
    char c;
    char* string;
    int size;
} StringParse;

StringParse parsers[] ={
        {'1', "one", 3},
        {'2', "two", 3},
        {'3', "three", 5},
        {'4', "four", 4},
        {'5', "five", 4},
        {'6', "six", 3},
        {'7', "seven", 5},
        {'8', "eight", 5},
        {'9', "nine", 4}
};


int main(){
    FILE* fp;
    char slice[5];
    char curr_chars[3] = {0, 0, '\0'};
    char curr_line[100];
    long int total = 0;
    int flag = 0;

    // 1. get file
    fp = fopen("../input.txt", "r");
    if(fp == NULL){
        perror("Error");
        exit(-1);
    }

    // iterate lines
    while(fgets(curr_line, sizeof(curr_line), fp) != NULL){
        // iterate line
        int i = 0;
        while(curr_line[i] != '\n' && i < sizeof(curr_line) && curr_line[i] != '\0'){
            memcpy(slice, &(curr_line[i]) , sizeof(slice));

            // check for digit
            if(slice[0] >= '0' && slice[0] <= '9'){
                if(!flag){
                    curr_chars[0] = slice[0];
                    flag = 1;
                } else curr_chars[1] = slice[0];
                i++;
                continue;
            }

            // check for spelled out digit
            for(int j = 0; j < sizeof(parsers) / sizeof(StringParse) ; j++){
                char* ptr = strstr(slice, parsers[j].string);
                if(ptr != NULL && ptr == &slice[0]){
                    if(!flag){
                        curr_chars[0] = parsers[j].c;
                        flag = 1;
                    } else curr_chars[1] = parsers[j].c;
                    i += parsers[j].size - 2;
                    break;
                }
            }
            i++;
        }
        if(curr_chars[1] == 0){
            curr_chars[1] = curr_chars[0];
        }
        flag = 0;
        total += atoi(curr_chars);
        curr_chars[0] = 0;
        curr_chars[1] = 0;
    }

    printf("Total : %ld\n", total);
}
