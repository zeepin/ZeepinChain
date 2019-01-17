# Wasm Smart CZPTract

## Introduction
Wasm (WebAssembly) is a binary instruction format for stack-based virtual machines. Wasm is designed to be a portable target for compiling high-level languages such as C/C++/Rust, and supports deployment of client and server applications on the web. ZPTology supports smart cZPTracts written in the Wasm format.

## Compilation

1. First, we will prepare a simple c language cZPTract that calculates the sum of two integers or concatenates two strings.

```c
//system apis
void * calloc(int count,int length);
void * malloc(int size);
int arrayLen(void *a);
int memcpy(void * dest,void * src,int length);
int memset(void * dest,char c,int length);

//utility apis
int strcmp(char *a,char *b);
char * fromcstring(char *s);
char * strconcat(char *a,char *b);
int Atoi(char * s);
long long Atoi64(char *s);
char * Itoa(int a);
char * I64toa(long long amount,int radix);
char * SHA1(char *s);
char * SHA256(char *s);

//parameter apis
int ZPT_ReadInt32Param(char *args);
long long ZPT_ReadInt64Param(char * args);
char * ZPT_ReadStringParam(char * args);
void ZPT_JsonUnmashalInput(void * addr,int size,char * arg);
char * ZPT_JsonMashalResult(void * val,char * types,int succeed);
char * ZPT_JsonMashalParams(void * s);
char * ZPT_RawMashalParams(void *s);
char * ZPT_GetCallerAddress();
char * ZPT_GetSelfAddress();
char * ZPT_MarshalNativeParams(void * s);
char * ZPT_MarshalNeoParams(void * s);

//Runtime apis
int ZPT_Runtime_CheckWitness(char * address);
void ZPT_Runtime_Notify(char * address);
int ZPT_Runtime_CheckSig(char * pubkey,char * data,char * sig);
int ZPT_Runtime_GetTime();
void ZPT_Runtime_Log(char * message);

//Attribute apis
int ZPT_Attribute_GetUsage(char * data);
char * ZPT_Attribute_GetData(char * data);

//Block apis
char * ZPT_Block_GetCurrentHeaderHash();
int ZPT_Block_GetCurrentHeaderHeight();
char * ZPT_Block_GetCurrentBlockHash();
int ZPT_Block_GetCurrentBlockHeight();
char * ZPT_Block_GetTransactionByHash(char * hash);
int * ZPT_Block_GetTransactionCountByBlkHash(char * hash);
int * ZPT_Block_GetTransactionCountByBlkHeight(int height);
char ** ZPT_Block_GetTransactionsByBlkHash(char * hash);
char ** ZPT_Block_GetTransactionsByBlkHeight(int height);


//Blockchain apis
int ZPT_BlockChain_GetHeight();
char * ZPT_BlockChain_GetHeaderByHeight(int height);
char * ZPT_BlockChain_GetHeaderByHash(char * hash);
char * ZPT_BlockChain_GetBlockByHeight(int height);
char * ZPT_BlockChain_GetBlockByHash(char * hash);
char * ZPT_BlockChain_GetCZPTract(char * address);

//header apis
char * ZPT_Header_GetHash(char * data);
int ZPT_Header_GetVersion(char * data);
char * ZPT_Header_GetPrevHash(char * data);
char * ZPT_Header_GetMerkleRoot(char  * data);
int ZPT_Header_GetIndex(char * data);
int ZPT_Header_GetTimestamp(char * data);
long long ZPT_Header_GetConsensusData(char * data);
char * ZPT_Header_GetNextConsensus(char * data);

//storage apis
void ZPT_Storage_Put(char * key,char * value);
char * ZPT_Storage_Get(char * key);
void ZPT_Storage_Delete(char * key);

//transaction apis
char * ZPT_Transaction_GetHash(char * data);
int ZPT_Transaction_GetType(char * data);
char * ZPT_Transaction_GetAttributes(char * data);


int add(int a, int b ){
        return a + b;
}

char * concat(char * a, char * b){
	int lena = arrayLen(a);
	int lenb = arrayLen(b);
	char * res = (char *)malloc((lena + lenb)*sizeof(char));
	for (int i = 0 ;i < lena ;i++){
		res[i] = a[i];
	}

	for (int j = 0; j < lenb ;j++){
		res[lena + j] = b[j];
	}
	return res;
}


int sumArray(int * a, int * b){

	int res = 0;
	int lena = arrayLen(a);
	int lenb = arrayLen(b);

	for (int i = 0;i<lena;i++){
		res += a[i];
	}
	for (int j = 0;j<lenb;j++){
		res += b[j];
	}
	return res;
}


char * invoke(char * method,char * args){

    if (strcmp(method ,"init")==0 ){
            return "init success!";
    }

    if (strcmp(method, "add")==0){
        int a = ZPT_ReadInt32Param(args);
	int b = ZPT_ReadInt32Param(args);
        int res = add(a,b);
        char * result = ZPT_JsonMashalResult(res,"int",1);
        ZPT_Runtime_Notify(result);
        return result;
    }

	if(strcmp(method,"concat")==0){
		char * a = ZPT_ReadStringParam(args);
		char * b = ZPT_ReadStringParam(args);
		char * res = concat(a,b);
		char * result = ZPT_JsonMashalResult(res,"string",1);
		ZPT_Runtime_Notify(result);
		return result;
	}
	if(strcmp(method,"addStorage")==0){
		char * a = ZPT_ReadStringParam(args);
		char * b = ZPT_ReadStringParam(args);
		ZPT_Storage_Put(a,b);
		char * result = ZPT_JsonMashalResult("Done","string",1);
		ZPT_Runtime_Notify(result);
		return result;
      }

      if(strcmp(method,"getStorage")==0){
		char * a = ZPT_ReadStringParam(args);
		char * value = ZPT_Storage_Get(a);
		char * result = ZPT_JsonMashalResult(value,"string",1);
		ZPT_Runtime_Notify(result);
		return result;
     }

     if(strcmp(method,"deleteStorage")==0){

        char * a = ZPT_ReadStringParam(args);
        ZPT_Storage_Delete(a);
        char * result = ZPT_JsonMashalResult("Done","string",1);
        ZPT_Runtime_Notify(result);
        return result;
    }
}


```

The following functions are provided by the virtual machine API and need to be declared at the head of the file.

```c
//system apis
void * calloc(int count,int length);
void * malloc(int size);
int arrayLen(void *a);
int memcpy(void * dest,void * src,int length);
int memset(void * dest,char c,int length);

//utility apis
int strcmp(char *a,char *b);
char * fromcstring(char *s);
char * strconcat(char *a,char *b);
int Atoi(char * s);
long long Atoi64(char *s);
char * Itoa(int a);
char * I64toa(long long amount,int radix);
char * SHA1(char *s);
char * SHA256(char *s);

//parameter apis
int ZPT_ReadInt32Param(char *args);
long long ZPT_ReadInt64Param(char * args);
char * ZPT_ReadStringParam(char * args);
void ZPT_JsonUnmashalInput(void * addr,int size,char * arg);
char * ZPT_JsonMashalResult(void * val,char * types,int succeed);
char * ZPT_JsonMashalParams(void * s);
char * ZPT_RawMashalParams(void *s);
char * ZPT_GetCallerAddress();
char * ZPT_GetSelfAddress();
char * ZPT_MarshalNativeParams(void * s);
char * ZPT_MarshalNeoParams(void * s);

//Runtime apis
int ZPT_Runtime_CheckWitness(char * address);
void ZPT_Runtime_Notify(char * address);
int ZPT_Runtime_CheckSig(char * pubkey,char * data,char * sig);
int ZPT_Runtime_GetTime();
void ZPT_Runtime_Log(char * message);

//Attribute apis
int ZPT_Attribute_GetUsage(char * data);
char * ZPT_Attribute_GetData(char * data);

//Block apis
char * ZPT_Block_GetCurrentHeaderHash();
int ZPT_Block_GetCurrentHeaderHeight();
char * ZPT_Block_GetCurrentBlockHash();
int ZPT_Block_GetCurrentBlockHeight();
char * ZPT_Block_GetTransactionByHash(char * hash);
int * ZPT_Block_GetTransactionCountByBlkHash(char * hash);
int * ZPT_Block_GetTransactionCountByBlkHeight(int height);
char ** ZPT_Block_GetTransactionsByBlkHash(char * hash);
char ** ZPT_Block_GetTransactionsByBlkHeight(int height);


//Blockchain apis
int ZPT_BlockChain_GetHeight();
char * ZPT_BlockChain_GetHeaderByHeight(int height);
char * ZPT_BlockChain_GetHeaderByHash(char * hash);
char * ZPT_BlockChain_GetBlockByHeight(int height);
char * ZPT_BlockChain_GetBlockByHash(char * hash);
char * ZPT_BlockChain_GetCZPTract(char * address);

//header apis
char * ZPT_Header_GetHash(char * data);
int ZPT_Header_GetVersion(char * data);
char * ZPT_Header_GetPrevHash(char * data);
char * ZPT_Header_GetMerkleRoot(char  * data);
int ZPT_Header_GetIndex(char * data);
int ZPT_Header_GetTimestamp(char * data);
long long ZPT_Header_GetConsensusData(char * data);
char * ZPT_Header_GetNextConsensus(char * data);

//storage apis
void ZPT_Storage_Put(char * key,char * value);
char * ZPT_Storage_Get(char * key);
void ZPT_Storage_Delete(char * key);

//transaction apis
char * ZPT_Transaction_GetHash(char * data);
int ZPT_Transaction_GetType(char * data);
char * ZPT_Transaction_GetAttributes(char * data);


```

The entry of Wasm contract is unified as ```char * invoke(char * method, char * args)```.

**method** is the method’s name that needs to be called.

**args** are the incoming parameters, raw bytes.

[More details on Wasm cZPTract APIs](wasm_api.md)

2. Compile the above C file into a smart contract in Wasm format.
    * Emscripten tool [http://kripken.github.io/emscripten-site/](http://http://kripken.github.io/emscripten-site/)
    * Online compiler WasmFiddle [https://wasdk.github.io/WasmFiddle](https://wasdk.github.io/WasmFiddle)
    * Use WasmFiddle as an Example
​    

Paste the C code into the edit window of "c". Please ignore the cZPTents of the "JS" window and click the "Build" button. If the compilation is correct, you can see the compiled Wasm format code in the "Text Format" window. If the compilation is wrong, an error message will be displayed in the "output" window.

![fiddle](images/fiddle.png)

If you are familiar with [wast syntax](http://webassembly.org/docs/binary-encoding/), you can modify the wast file yourself.

And use the [wabt](https://github.com/WebAssembly/wabt) tool to compile the wast file into Wasm format.

​    
3. Click the "Wasm" button to download the compiled Wasm file.


### Passing parameters in JSON format

Incoming parameters can be in JSON format:

```c
//system apis
void * calloc(int count,int length);
void * malloc(int size);
int arrayLen(void *a);
int memcpy(void * dest,void * src,int length);
int memset(void * dest,char c,int length);

//utility apis
int strcmp(char *a,char *b);
char * fromcstring(char *s);
char * strconcat(char *a,char *b);
int Atoi(char * s);
long long Atoi64(char *s);
char * Itoa(int a);
char * I64toa(long long amount,int radix);
char * SHA1(char *s);
char * SHA256(char *s);

//parameter apis
int ZPT_ReadInt32Param(char *args);
long long ZPT_ReadInt64Param(char * args);
char * ZPT_ReadStringParam(char * args);
void ZPT_JsonUnmashalInput(void * addr,int size,char * arg);
char * ZPT_JsonMashalResult(void * val,char * types,int succeed);
char * ZPT_JsonMashalParams(void * s);
char * ZPT_RawMashalParams(void *s);
char * ZPT_GetCallerAddress();
char * ZPT_GetSelfAddress();
char * ZPT_MarshalNativeParams(void * s);
char * ZPT_MarshalNeoParams(void * s);

//Runtime apis
int ZPT_Runtime_CheckWitness(char * address);
void ZPT_Runtime_Notify(char * address);
int ZPT_Runtime_CheckSig(char * pubkey,char * data,char * sig);
int ZPT_Runtime_GetTime();
void ZPT_Runtime_Log(char * message);

//Attribute apis
int ZPT_Attribute_GetUsage(char * data);
char * ZPT_Attribute_GetData(char * data);

//Block apis
char * ZPT_Block_GetCurrentHeaderHash();
int ZPT_Block_GetCurrentHeaderHeight();
char * ZPT_Block_GetCurrentBlockHash();
int ZPT_Block_GetCurrentBlockHeight();
char * ZPT_Block_GetTransactionByHash(char * hash);
int * ZPT_Block_GetTransactionCountByBlkHash(char * hash);
int * ZPT_Block_GetTransactionCountByBlkHeight(int height);
char ** ZPT_Block_GetTransactionsByBlkHash(char * hash);
char ** ZPT_Block_GetTransactionsByBlkHeight(int height);


//Blockchain apis
int ZPT_BlockChain_GetHeight();
char * ZPT_BlockChain_GetHeaderByHeight(int height);
char * ZPT_BlockChain_GetHeaderByHash(char * hash);
char * ZPT_BlockChain_GetBlockByHeight(int height);
char * ZPT_BlockChain_GetBlockByHash(char * hash);
char * ZPT_BlockChain_GetCZPTract(char * address);

//header apis
char * ZPT_Header_GetHash(char * data);
int ZPT_Header_GetVersion(char * data);
char * ZPT_Header_GetPrevHash(char * data);
char * ZPT_Header_GetMerkleRoot(char  * data);
int ZPT_Header_GetIndex(char * data);
int ZPT_Header_GetTimestamp(char * data);
long long ZPT_Header_GetConsensusData(char * data);
char * ZPT_Header_GetNextConsensus(char * data);

//storage apis
void ZPT_Storage_Put(char * key,char * value);
char * ZPT_Storage_Get(char * key);
void ZPT_Storage_Delete(char * key);

//transaction apis
char * ZPT_Transaction_GetHash(char * data);
int ZPT_Transaction_GetType(char * data);
char * ZPT_Transaction_GetAttributes(char * data);
int add(int a, int b ){
        return a + b;
}

char * concat(char * a, char * b){
	int lena = arrayLen(a);
	int lenb = arrayLen(b);
	char * res = (char *)malloc((lena + lenb)*sizeof(char));
	for (int i = 0 ;i < lena ;i++){
		res[i] = a[i];
	}

	for (int j = 0; j < lenb ;j++){
		res[lena + j] = b[j];
	}
	return res;
}


int sumArray(int * a, int * b){

	int res = 0;
	int lena = arrayLen(a);
	int lenb = arrayLen(b);

	for (int i = 0;i<lena;i++){
		res += a[i];
	}
	for (int j = 0;j<lenb;j++){
		res += b[j];
	}
	return res;
}


char * invoke(char * method,char * args){

    if (strcmp(method ,"init")==0 ){
            return "init success!";
    }

    if (strcmp(method, "add")==0){
        struct Params {
                int a;
                int b;
        };
        struct Params *p = (struct Params *)malloc(sizeof(struct Params));

        ZPT_JsonUnmashalInput(p,sizeof(struct Params),args);
        int res = add(p->a,p->b);
        char * result = ZPT_JsonMashalResult(res,"int",1);
        ZPT_Runtime_Notify(result);
        return result;
    }

	if(strcmp(method,"concat")==0){
		struct Params{
			char *a;
			char *b;
		};
		struct Params *p = (struct Params *)malloc(sizeof(struct Params));
		ZPT_JsonUnmashalInput(p,sizeof(struct Params),args);
		char * res = concat(p->a,p->b);
		char * result = ZPT_JsonMashalResult(res,"string",1);
		ZPT_Runtime_Notify(result);
		return result;
	}
	
	if(strcmp(method,"sumArray")==0){
		struct Params{
			int *a;
			int *b;
		};
		struct Params *p = (struct Params *)malloc(sizeof(struct Params));
		ZPT_JsonUnmashalInput(p,sizeof(struct Params),args);
		int res = sumArray(p->a,p->b);
		char * result = ZPT_JsonMashalResult(res,"int",1);
		ZPT_Runtime_Notify(result);
		return result;
	}

	if(strcmp(method,"addStorage")==0){

		struct Params{
			char * a;
			char * b;
		};
		struct Params *p = (struct Params *)malloc(sizeof(struct Params));
		ZPT_JsonUnmashalInput(p,sizeof(struct Params),args);
		ZPT_Storage_Put(p->a,p->b);
		char * result = ZPT_JsonMashalResult("Done","string",1);
		ZPT_Runtime_Notify(result);
		return result;
    }

	if(strcmp(method,"getStorage")==0){
		struct Params{
			char * a;
		};
		struct Params *p = (struct Params *)malloc(sizeof(struct Params));
		ZPT_JsonUnmashalInput(p,sizeof(struct Params),args);
		char * value = ZPT_Storage_Get(p->a);
		char * result = ZPT_JsonMashalResult(value,"string",1);
		ZPT_Runtime_Notify(result);
		return result;
	}

	if(strcmp(method,"deleteStorage")==0){

        struct Params{
                char * a;
        };
		struct Params *p = (struct Params *)malloc(sizeof(struct Params));
		ZPT_JsonUnmashalInput(p,sizeof(struct Params),args);
        ZPT_Storage_Delete(p->a);
        char * result = ZPT_JsonMashalResult("Done","string",1);
        ZPT_Runtime_Notify(result);
        return result;
    }
}
                                                                                                                                      
```



### Handle Blockchain functi0on


```c
//system apis
void * calloc(int count,int length);
void * malloc(int size);
int arrayLen(void *a);
int memcpy(void * dest,void * src,int length);
int memset(void * dest,char c,int length);

//utility apis
int strcmp(char *a,char *b);
char * fromcstring(char *s);
char * strconcat(char *a,char *b);
int Atoi(char * s);
long long Atoi64(char *s);
char * Itoa(int a);
char * I64toa(long long amount,int radix);
char * SHA1(char *s);
char * SHA256(char *s);

//parameter apis
int ZPT_ReadInt32Param(char *args);
long long ZPT_ReadInt64Param(char * args);
char * ZPT_ReadStringParam(char * args);
void ZPT_JsonUnmashalInput(void * addr,int size,char * arg);
char * ZPT_JsonMashalResult(void * val,char * types,int succeed);
char * ZPT_JsonMashalParams(void * s);
char * ZPT_RawMashalParams(void *s);
char * ZPT_GetCallerAddress();
char * ZPT_GetSelfAddress();
char * ZPT_MarshalNativeParams(void * s);
char * ZPT_MarshalNeoParams(void * s);

//Runtime apis
int ZPT_Runtime_CheckWitness(char * address);
void ZPT_Runtime_Notify(char * address);
int ZPT_Runtime_CheckSig(char * pubkey,char * data,char * sig);
int ZPT_Runtime_GetTime();
void ZPT_Runtime_Log(char * message);

//Attribute apis
int ZPT_Attribute_GetUsage(char * data);
char * ZPT_Attribute_GetData(char * data);

//Block apis
char * ZPT_Block_GetCurrentHeaderHash();
int ZPT_Block_GetCurrentHeaderHeight();
char * ZPT_Block_GetCurrentBlockHash();
int ZPT_Block_GetCurrentBlockHeight();
char * ZPT_Block_GetTransactionByHash(char * hash);
int * ZPT_Block_GetTransactionCountByBlkHash(char * hash);
int * ZPT_Block_GetTransactionCountByBlkHeight(int height);
char ** ZPT_Block_GetTransactionsByBlkHash(char * hash);
char ** ZPT_Block_GetTransactionsByBlkHeight(int height);


//Blockchain apis
int ZPT_BlockChain_GetHeight();
char * ZPT_BlockChain_GetHeaderByHeight(int height);
char * ZPT_BlockChain_GetHeaderByHash(char * hash);
char * ZPT_BlockChain_GetBlockByHeight(int height);
char * ZPT_BlockChain_GetBlockByHash(char * hash);
char * ZPT_BlockChain_GetCZPTract(char * address);

//header apis
char * ZPT_Header_GetHash(char * data);
int ZPT_Header_GetVersion(char * data);
char * ZPT_Header_GetPrevHash(char * data);
char * ZPT_Header_GetMerkleRoot(char  * data);
int ZPT_Header_GetIndex(char * data);
int ZPT_Header_GetTimestamp(char * data);
long long ZPT_Header_GetConsensusData(char * data);
char * ZPT_Header_GetNextConsensus(char * data);

//storage apis
void ZPT_Storage_Put(char * key,char * value);
char * ZPT_Storage_Get(char * key);
void ZPT_Storage_Delete(char * key);

//transaction apis
char * ZPT_Transaction_GetHash(char * data);
int ZPT_Transaction_GetType(char * data);
char * ZPT_Transaction_GetAttributes(char * data);
int add(int a, int b ){
        return a + b;
}

char * concat(char * a, char * b){
	int lena = arrayLen(a);
	int lenb = arrayLen(b);
	char * res = (char *)malloc((lena + lenb)*sizeof(char));
	for (int i = 0 ;i < lena ;i++){
		res[i] = a[i];
	}

	for (int j = 0; j < lenb ;j++){
		res[lena + j] = b[j];
	}
	return res;
}


int sumArray(int * a, int * b){

	int res = 0;
	int lena = arrayLen(a);
	int lenb = arrayLen(b);

	for (int i = 0;i<lena;i++){
		res += a[i];
	}
	for (int j = 0;j<lenb;j++){
		res += b[j];
	}
	return res;
}


char * invoke(char * method,char * args){

    if (strcmp(method ,"init")==0 ){
            return "init success!";
    }

    if (strcmp(method, "getBlockHeaderHash")==0){
        struct Params {
                int height;
                char * blockhash;
        };
        struct Params *p = (struct Params *)malloc(sizeof(struct Params));

        ZPT_JsonUnmashalInput(p,sizeof(struct Params),args);
        char *data=ZPT_BlockChain_GetHeaderByHeight(p->height);
        char * hash=ZPT_Header_GetHash(data);
        char * result = ZPT_JsonMashalResult(hash,"string",1);
        ZPT_Runtime_Notify(result);
        return result;
    }
    if (strcmp(method, "getTransactionHash")==0){
        struct Params {
                int height;
                char * transactionhash;
        };
        struct Params *p = (struct Params *)malloc(sizeof(struct Params));

        ZPT_JsonUnmashalInput(p,sizeof(struct Params),args);
        char* transactionhash=fromcstring(p->transactionhash);
        char *data=ZPT_Block_GetTransactionByHash(transactionhash);
        char * hash=ZPT_Transaction_GetHash(data);
        char * result = ZPT_JsonMashalResult(hash,"string",1);
        ZPT_Runtime_Notify(result);
        return result;
    }
}
                                                                                                                                      
```
